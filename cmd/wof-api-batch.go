package main

// https://github.com/whosonfirst/whosonfirst-www-api/issues/99#issuecomment-333960724

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/client"
	"github.com/whosonfirst/go-whosonfirst-api/endpoint"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type BatchRequestSet struct {
	APIKey   string
	Requests []BatchRequest
}

type BatchRequest map[string]string

type BatchResponse struct {
	Index       int
	APIResponse api.APIResponse
	Timing      time.Duration
}

// please for to be passing in a timeout context here and to make
// sure it bubbles down to any individual requests being processed

func process_batch(rs BatchRequestSet) ([]api.APIResponse, error) {

	e, err := endpoint.NewMapzenAPIEndpoint(rs.APIKey)

	if err != nil {
		return nil, err
	}

	c, _ := client.NewHTTPClient(e)

	response_ch := make(chan BatchResponse)
	error_ch := make(chan error)
	done_ch := make(chan bool)

	complete_ch := make(chan bool)

	pending := len(rs.Requests)
	responses := make([]api.APIResponse, pending)

	go func() {

		for pending > 0 {

			select {
			case rsp := <-response_ch:
				responses[rsp.Index] = rsp.APIResponse
			case err := <-error_ch:
				log.Println(err)
			case <-done_ch:
				pending -= 1
			}
		}

		complete_ch <- true
	}()

	t1 := time.Now()
	
	for idx, req := range rs.Requests {

		// please for to be rate-limiting here...

		go process_request(c, idx, req, response_ch, error_ch, done_ch)
	}
	
	<- complete_ch

	t2 := time.Since(t1)

	log.Println("TIME", t2)
	log.Println("RESPONSES", responses)
	
	return responses, nil
}

func process_request(c api.APIClient, idx int, req BatchRequest, response_ch chan BatchResponse, error_ch chan error, done_ch chan bool) {

	defer func() {
		done_ch <- true
	}()

	t1 := time.Now()

	cb := func(rsp api.APIResponse) error {

		t2 := time.Since(t1)

		response := BatchResponse{
			Index:       idx,
			APIResponse: rsp,
			Timing:      t2,
		}

		response_ch <- response
		return nil
	}

	method := ""
	args := url.Values{}

	for k, v := range req {

		if k == "method" {
			method = v
			continue
		}

		args.Set(k, v)
	}

	if method == "" {
		error_ch <- errors.New("Missing API method")
		return
	}

	err := c.ExecuteMethodWithCallback(method, &args, cb)

	if err != nil {
		error_ch <- err
	}
}

func parse(body []byte) ([]BatchRequest, error) {

	batch := make([]BatchRequest, 0)

	c := gjson.GetBytes(body, "#")
	count := int(c.Int())

	if count == 0 {
		return nil, errors.New("Invalid batch request")
	}

	for i := 0; i < count; i++ {

		path := strconv.Itoa(i)
		r := gjson.GetBytes(body, path)

		br := make(map[string]string)

		for k, v := range r.Map() {
			br[k] = v.String()
		}

		batch = append(batch, br)
	}

	if len(batch) == 0 {
		return nil, errors.New("Invalid batch request")
	}

	return batch, nil
}

func handler() (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		query := req.URL.Query()

		api_key := query.Get("api_key")
		api_key = strings.Trim(api_key, " ")

		if api_key == "" {
			http.Error(rsp, "Missing API key", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		// hash body here and ensure that we don't already have a running
		// batch for api_key + "#" + hash

		requests, err := parse(body)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		request_set := BatchRequestSet{
			APIKey:   api_key,
			Requests: requests,
		}

		// log api_key + "#" + hash here - it would be nice to all of this using
		// BatchRequestSet but that means always parsing body first...
		
		// see notes above wrt a timeout context (as in: it does not exist yet)

		res, err := process_batch(request_set)

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		js, err := json.Marshal(res)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.Write(js)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	flag.Parse()

	// fetch list of valid API methods from the api.spec.method and pass along
	// to handler for basic validation on all requests here...
	
	batch_handler, err := handler()

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", batch_handler)

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("listening on %s\n", endpoint)

	err = http.ListenAndServe(endpoint, mux)

	if err != nil {
		log.Fatal(err)
	}

}
