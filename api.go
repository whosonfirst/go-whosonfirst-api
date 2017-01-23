package api

import (
	_ "fmt"
	"github.com/jeffail/gabs"
	"io/ioutil"
	"net/http"
	"net/url"
)

type APIClient struct {
	scheme   string
	isa      string
	Host     string
	Endpoint string
	Token    string
}

type APIError struct {
	Code    int64
	Message string
}

type APIResponse struct {
	Raw    []byte
	Parsed *gabs.Container
}

func (rsp APIResponse) Stat() string {

	var v string

	v, _ = rsp.Parsed.Path("stat").Data().(string)
	return v
}

func (rsp APIResponse) Ok() (bool, *APIError) {

	stat := rsp.Stat()

	if stat == "ok" {
		return true, nil
	}

	return false, rsp.Error()
}

func (rsp APIResponse) Error() *APIError {

	var code int64
	var msg string

	// why does this (lookup for error.code) always return 0?

	code, _ = rsp.Parsed.Path("error.code").Data().(int64)
	msg, _ = rsp.Parsed.Path("error.message").Data().(string)

	err := APIError{Code: code, Message: msg}
	return &err
}

func (rsp APIResponse) Body() *gabs.Container {
	return rsp.Parsed
}

func (rsp APIResponse) Dumps() string {
	return rsp.Parsed.String()
}

func ParseAPIResponse(raw []byte) (*APIResponse, error) {

	parsed, parse_err := gabs.ParseJSON(raw)

	if parse_err != nil {
		return nil, parse_err
	}

	rsp := APIResponse{
		Raw:    raw,
		Parsed: parsed,
	}

	return &rsp, nil
}

func OAuth2Client(token string) *APIClient {

	return &APIClient{
		isa:      "oauth2",
		scheme:   "https",
		Token:    token,
		Host:     "",
		Endpoint: "",
	}
}

func MapzenAPIClient(key string) *APIClient {

	return &APIClient{
		isa:      "mapzen-api",
		scheme:   "https",
		Token:    key,
		Host:     "whosonfirst-api.mapzen.com",
		Endpoint: "",
	}
}

func (client *APIClient) ExecuteMethod(method string, params *url.Values) (*APIResponse, error) {

	url := client.scheme + "://" + client.Host + "/" + client.Endpoint

	params.Set("method", method)

	if client.isa == "mapzen-api" {
		params.Set("api_key", client.Token)
	} else {
		params.Set("access_token", client.Token)
	}
	
	http_req, req_err := http.NewRequest("POST", url, nil)

	if req_err != nil {
		return nil, req_err
	}

	http_req.URL.RawQuery = (*params).Encode()

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
