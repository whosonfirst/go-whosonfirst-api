package api

import (
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
	Stat() string
	Ok() (bool, APIError)
	// Get(string) interface{}
}

type APIError interface {
	String() string
	Code() int64
	Message() string
}
