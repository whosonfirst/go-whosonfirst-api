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
	Ok() (bool, APIError)
	Pagination() (APIPagination, error)
}

type APIError interface {
	String() string
	Code() int64
	Message() string
}

type APIPagination interface {
	String() string
	Pages() int
	Page() int
	PerPage() int
	Total() int
}

type APIResponseCallback func(APIResponse) error
