package client

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/response"
	_ "log"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient struct {
	api.APIClient
	endpoint api.APIEndpoint
	qpslimit int
	qpmlimit int
	qphlimit int
	qpscount int
	qpmcount int
	qphcount int
}

func NewHTTPClient(endpoint api.APIEndpoint) (*HTTPClient, error) {

	cl := HTTPClient{
		endpoint: endpoint,
		qpslimit: 6,
		qpmlimit: 30,
		qphlimit: 1000,
		qpscount: 0,
		qpmcount: 0,
		qphcount: 0,
	}

	// TO DO: set up a channel/throttle to block calls to ExecuteMethod
	// from exceeding QPS/M/H (20170125/thisisaaronland)

	return &cl, nil
}

func (client *HTTPClient) DefaultArgs() *url.Values {
	args := url.Values{}
	return &args
}

func (client *HTTPClient) ExecuteMethod(method string, params *url.Values) (api.APIResponse, error) {

	params.Set("method", method)

	http_req, err := client.endpoint.NewRequest(params)

	if err != nil {
		return nil, err
	}

	http_req.Header.Add("Accept-Encoding", "gzip")

	http_client := &http.Client{}
	http_rsp, http_err := http_client.Do(http_req)

	if http_err != nil {
		return nil, http_err
	}

	defer http_rsp.Body.Close()

	switch http_rsp.StatusCode {

	case 200:
		// pass
	case 201:
		// pass
	default:

		header := http_rsp.Header
		msg := fmt.Sprintf("%s %s", header.Get("X-Api-Error-Code"), header.Get("X-Api-Error-Message"))

		return nil, errors.New(msg)
	}

	var rsp api.APIResponse
	var parse_err error

	switch params.Get("format") {

	case "":
		rsp, parse_err = response.ParseJSONResponse(http_rsp)
	case "csv":
		rsp, parse_err = response.ParseCSVResponse(http_rsp)
	case "json":
		rsp, parse_err = response.ParseJSONResponse(http_rsp)
	case "meta":
		rsp, parse_err = response.ParseMetaResponse(http_rsp)
	default:
		return nil, errors.New("Unsupported output format")
	}

	if parse_err != nil {
		return nil, parse_err
	}

	return rsp, nil
}

func (client *HTTPClient) ExecuteMethodWithCallback(method string, params *url.Values, callback api.APIResponseCallback) error {

	rsp, err := client.ExecuteMethod(method, params)

	if err != nil {
		return err
	}

	_, api_err := rsp.Ok()

	if api_err != nil {
		return errors.New(api_err.String())
	}

	return callback(rsp)
}

func (client *HTTPClient) ExecuteMethodPaginated(method string, params *url.Values, callback api.APIResponseCallback) error {

	api_key := params.Get("api_key") // PLEASE MAKE ME GENERIC AND INTERFACE-Y

	for {

		rsp, err := client.ExecuteMethod(method, params)

		if err != nil {
			return err
		}

		_, api_err := rsp.Ok()

		if api_err != nil {
			return errors.New(api_err.String())
		}

		pg, err := rsp.Pagination()

		if err != nil {
			return err
		}

		next_query := pg.NextQuery()

		cb_err := callback(rsp)

		if cb_err != nil {
			return cb_err
		}

		if next_query == "null" || next_query == "" {
			break
		}

		parsed, err := url.ParseQuery(next_query)

		if err != nil {
			return err
		}

		parsed.Set("api_key", api_key) // SEE ABOVE ABOUT GENERIC AND INTERFACE-Y

		params = &parsed

		// to do: add proper QPS throttling here

		time.Sleep(200 * time.Millisecond)
	}

	return nil
}
