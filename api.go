package api

import (
	"github.com/tidwall/gjson" // see notes below
	"net/http"
	"net/url"
)

type APIEndpoint interface {
	URL() (*url.URL, error)
	NewRequest(*url.Values) (*http.Request, error)
}

type APIClient interface {
	ExecuteMethod(string, *url.Values) (APIResponse, error)
	DefaultArgs() *url.Values
}

type APIResponse interface {
	String() string
	Ok() (bool, APIError)
	Get(string) gjson.Result // this is not correct - we need something more generic but interface{} doesn't work...
	// we need something something something pagination here...
}

type APIError interface {
	String() string
	Code() int64
	Message() string
}
