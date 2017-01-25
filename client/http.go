package client

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/response"
	"io/ioutil"
	_ "log"
	"net/http"
	"net/url"
	"strconv"
	_ "sync/atomic"
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

	format := params.Get("format")

	if format != "" && format != "json" {
		return nil, errors.New("JSON is the only output format currently supported")
	}

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

	http_body, io_err := ioutil.ReadAll(http_rsp.Body)

	if io_err != nil {
		return nil, io_err
	}

	// to do: support other formats...

	rsp, parse_err := response.ParseJSONResponse(http_body)

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

	pages := 0
	page := 1

	for {

		params.Set("page", strconv.Itoa(page))

		rsp, err := client.ExecuteMethod(method, params)

		if err != nil {
			return err
		}

		_, api_err := rsp.Ok()

		if api_err != nil {
			return errors.New(api_err.String())
		}

		if pages == 0 {

			pg, err := rsp.Pagination()

			if err != nil {
				return err
			}

			pages = pg.Pages()
		}

		cb_err := callback(rsp)

		if cb_err != nil {
			return cb_err
		}

		if page >= pages {
			break
		}

		// to do: add proper QPS throttling here

		time.Sleep(200 * time.Millisecond)

		page += 1
	}

	return nil
}
