package endpoint

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-api"
	"net/http"
	"net/url"
)

type MapzenAPIEndpoint struct {
	api.APIEndpoint
	Scheme string
	Host   string
	Path   string
	Key    string
}

func NewMapzenAPIEndpoint(key string) (*MapzenAPIEndpoint, error) {

	e := MapzenAPIEndpoint{
		Scheme: "https",
		Host:   "places.mapzen.com",
		Path:   "v1",
		Key:    key,
	}

	return &e, nil
}

func (e *MapzenAPIEndpoint) SetEndpoint(custom string) error {

	parsed, err := url.Parse(custom)

	if err != nil {
		return err
	}

	e.Scheme = parsed.Scheme
	e.Host = parsed.Host
	e.Path = parsed.Path

	return nil
}

func (e *MapzenAPIEndpoint) URL() (*url.URL, error) {

	raw := fmt.Sprintf("%s://%s/%s", e.Scheme, e.Host, e.Path)

	url, err := url.Parse(raw)

	if err != nil {
		return nil, err
	}

	return url, err
}

func (e *MapzenAPIEndpoint) NewRequest(params *url.Values) (*http.Request, error) {

	url, err := e.URL()

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url.String(), nil)

	if err != nil {
		return nil, err
	}

	params.Set("api_key", e.Key)
	req.URL.RawQuery = (*params).Encode()

	return req, nil
}
