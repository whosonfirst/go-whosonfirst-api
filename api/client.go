package api

import (
       "io/ioutil"
       "net/http"
       "net/url"
)

type APIClient struct {
     endpoint APIEndpoint
}

func NewAPIClient (endpoint APIEndpoint) (*APIClient, error) {

     cl := APIClient{
     	endpoint: endpoint,
     }

     return &cl, nil
}

func (client *APIClient) ExecuteMethod(method string, params *url.Values) (*APIResponse, error) {

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

	rsp, parse_err := ParseAPIResponse(http_body)

	if parse_err != nil {
		return nil, parse_err
	}

	return rsp, nil
}
