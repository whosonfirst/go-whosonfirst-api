package api

import (
	"net/http"
	"net/url"
	"strings"
)

type APIEndpoint interface {
	URL() (*url.URL, error)
	NewRequest(*url.Values) (*http.Request, error)
	SetEndpoint(string) error
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

type APIPagination interface {
	String() string
	Pages() int
	Page() int
	PerPage() int
	Total() int
	Cursor() string
	NextQuery() string
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

type APIResultMultiWriter interface { // PLEASE RENAME ME...
	Write(APIResult) (int, error)
	Close()
}

type APIResultWriter interface {
	Write([]byte) (int, error)
	WriteString(string) (int, error)
	WriteResult(APIResult) (int, error)
	Close() error
}

type APIResponseCallback func(APIResponse) error

type APIParams []string

func (p *APIParams) String() string {
	return strings.Join(*p, "\n")
}

func (p *APIParams) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func (p *APIParams) ToArgs() *url.Values {

	args := url.Values{}

	for _, str_pair := range *p {
		pair := strings.Split(str_pair, "=")
		args.Set(pair[0], pair[1])
	}

	return &args
}
