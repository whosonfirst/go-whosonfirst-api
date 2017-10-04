package main

// https://github.com/whosonfirst/whosonfirst-www-api/issues/99#issuecomment-333960724

import (
	"encoding/json"
	"errors"
	_ "flag"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/client"
	"github.com/whosonfirst/go-whosonfirst-api/endpoint"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Request map[string]string

type BatchRequest []Request

type Response struct {
	Index       int
	APIResponse api.APIResponse
}

func process_batch(batch *BatchRequest) ([]api.APIResponse, error) {

	api_key := "..."

	e, err := endpoint.NewMapzenAPIEndpoint(api_key)

	if err != nil {
		return nil, err
	}

	c, _ := client.NewHTTPClient(e)

	response_ch := make(chan Response)
	error_ch := make(chan error)
	done_ch := make(chan bool)

	for idx, req := range batch {
		go process_request(c, idx, req, response_ch, error_ch, done_ch)
	}

	responses := make([]api.APIResponse, 0)
	pending := len(batch)

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

	return responses, nil
}

func process_request(c api.APIClient, idx int, req Request, response_ch chan Response, error_ch chan error, done_ch chan bool) {

	defer func() {
		done_ch <- true
	}()

	cb := func(rsp api.APIResponse) error {

		response := Response{
			Index:       idx,
			APIResponse: rsp,
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

func handler() (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		query := req.URL.Query()

		api_key := query.Get("api_key")
		api_key = strings.Trim(api_key, " ")

		if api_key == "" {
			http.Error(rsp, "Missing API key", http.StatusBadRequest)
			return
		}

		batch := query.Get("batch")
		batch = strings.Trim(batch, " ")

		if batch == "" {
			http.Error(rsp, "batch", http.StatusBadRequest)
			return
		}

		var batch_req BatchRequest
		err := json.Unmarshal([]byte(batch), &batch_req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := process(&batch_req)

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		// rsp.Write(enc)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}

func main() {

	batch_handler, err := handler()
}
