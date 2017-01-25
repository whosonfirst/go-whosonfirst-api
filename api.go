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
	Raw() []byte
	String() string
	Ok() (bool, APIError)
	Pagination() (APIPagination, error)
	Results() ([]APIResult, error)
}

type APIError interface {
	String() string
	Code() int64
	Message() string
}

type APIResult interface {
	WOFId() int64
	WOFParentId() int64
	WOFName() string
	WOFPlacetype() string
	WOFCountry() string
	WOFRepo() string
	Path() string
	URI() string
}

type APIPagination interface {
	String() string
	Pages() int
	Page() int
	PerPage() int
	Total() int
}

type APIResultWriter interface {
	Write([]byte) (int, error)
	WriteString(string) (int, error)
	WriteResult(APIResult) (int, error)
}

type APIResponseCallback func(APIResponse) error
