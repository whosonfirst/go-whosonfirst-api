package api

import (
	"fmt"
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
	ExecuteMethodWithCallback(string, *url.Values, APIResponseCallback) error
	DefaultArgs() *url.Values
}

type APIResponse interface {
	Raw() []byte
	String() string
	Ok() (bool, APIError)
	Pagination() (APIPagination, error)
	Places() ([]APIPlacesResult, error)
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

type APIPlacesResult interface {
	WOFId() int64
	WOFParentId() int64
	WOFName() string
	WOFPlacetype() string
	WOFCountry() string
	WOFRepo() string
	Path() string
	URI() string
	String(...APIResultFlag) string
}

// APIResultFlag is all wet paint... still trying to work it out

type APIResultFlag interface {
	Key() string
	Value() interface{}
	Bool() bool
	String() string
	Int() int
	Int64() int64
	Float64() float64
}

type APIResultBooleanFlag struct {
	APIResultFlag
	flkey   string
	flvalue bool
}

func (fl APIResultBooleanFlag) Key() string {
	return fl.flkey
}

func (fl APIResultBooleanFlag) Value() interface{} {
	return fl.flvalue
}

func (fl APIResultBooleanFlag) Bool() bool {
	return fl.Value().(bool)
}

func (fl APIResultBooleanFlag) String() string {
	return fmt.Sprintf("%t", fl.Bool())
}

func (fl APIResultBooleanFlag) Int() int {

	if fl.Bool() {
		return 1
	}

	return 0
}

func (fl APIResultBooleanFlag) Int64() int64 {
	return int64(fl.Int())
}

func (fl APIResultBooleanFlag) Float64() float64 {
	return float64(fl.Int())
}

func NewAPIResultBooleanFlag(key string, flag bool) APIResultBooleanFlag {

	fl := APIResultBooleanFlag{
		flkey:   key,
		flvalue: flag,
	}

	return fl
}

type APIResultMultiWriter interface { // PLEASE RENAME ME...
	Write(APIPlacesResult) (int, error)
	Close()
}

type APIResultWriter interface {
	Write([]byte) (int, error)
	WriteString(string) (int, error)
	WriteResult(APIPlacesResult) (int, error)
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
