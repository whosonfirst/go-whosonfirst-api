package client

import (
       "github.com/whosonfirst/go-whosonfirst-api"
       "github.com/whosonfirst/go-whosonfirst-api/response"       
	"io/ioutil"
	"net/http"
	"net/url"
)

type HTTPClient struct {
        api.APIClient     		
	endpoint api.APIEndpoint
}

func NewHTTPClient(endpoint api.APIEndpoint) (*HTTPClient, error) {

	cl := HTTPClient{
		endpoint: endpoint,
	}

	return &cl, nil
}

func (client *HTTPClient) DefaultArgs() *url.Values {
	args := url.Values{}
	return &args
}

func (client *HTTPClient) ExecuteMethod(method string, params *url.Values) (*response.APIResponse, error) {

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

	rsp, parse_err := response.ParseAPIResponse(http_body)

	if parse_err != nil {
		return nil, parse_err
	}

	return rsp, nil
}
