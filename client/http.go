package client

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/mapzen"
	"github.com/whosonfirst/go-whosonfirst-api/response"
	_ "log"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient struct {
	api.APIClient
	endpoint    api.APIEndpoint
	qpslimit    int
	qpmlimit    int
	qphlimit    int
	qpscount    int
	qpmcount    int
	qphcount    int
	http_client *http.Client
}

func NewHTTPClient(endpoint api.APIEndpoint) (*HTTPClient, error) {

	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

	http_client := &http.Client{Transport: tr}

	cl := HTTPClient{
		endpoint:    endpoint,
		qpslimit:    6,
		qpmlimit:    30,
		qphlimit:    1000,
		qpscount:    0,
		qpmcount:    0,
		qphcount:    0,
		http_client: http_client,
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

	http_rsp, http_err := client.http_client.Do(http_req)

	if http_err != nil {
		msg := fmt.Sprintf("HTTP request failed: %s", http_err.Error())
		return nil, errors.New(msg)
	}

	defer http_rsp.Body.Close()

	status_code := http_rsp.StatusCode

	if IsHTTPError(status_code) {

		if IsWOFError(status_code) {

			header := http_rsp.Header
			msg := fmt.Sprintf("%s %s", header.Get("X-Api-Error-Code"), header.Get("X-Api-Error-Message"))

			return nil, errors.New(msg)
		}

		if IsMapzenError(status_code) {

			rsp, err := mapzen.ParseMapzenResponse(http_rsp)

			if err != nil {
				return nil, err
			}

			msg := fmt.Sprintf("%d %s", rsp.Meta.StatusCode, rsp.Results.Error.Message)
			return nil, errors.New(msg)
		}

		return nil, errors.New(http_rsp.Status)
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

func IsHTTPError(status_code int) bool {
	return (status_code > 400 && status_code <= 599)
}

func IsMapzenError(status_code int) bool {

	if status_code >= 400 && status_code <= 404 {
		return true
	}

	if status_code == 429 {
		return true
	}

	if status_code == 502 {
		return true
	}

	if status_code == 504 {
		return true
	}

	return false
}

func IsWOFError(status_code int) bool {

	if status_code == 450 {
		return true
	}

	if status_code >= 452 && status_code <= 499 {
		return true
	}

	if status_code == 512 {
		return true
	}

	return false
}
