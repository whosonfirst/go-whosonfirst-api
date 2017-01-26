package endpoint

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-api"
	"net/http"
	"net/url"
)

type OAuth2APIEndpoint struct {
	api.APIEndpoint
	Scheme string
	Host   string
	Path   string
	Token  string
}

func NewOAuth2APIEndpoint(access_token string) (*OAuth2APIEndpoint, error) {

	e := OAuth2APIEndpoint{
		Scheme: "https",
		Host:   "",
		Path:   "",
		Token:  access_token,
	}

	return &e, nil
}

func (e *OAuth2APIEndpoint) SetEndpoint(custom string) error {

	parsed, err := url.Parse(custom)

	if err != nil {
		return err
	}

	e.Scheme = parsed.Scheme
	e.Host = parsed.Host
	e.Path = parsed.Path

	return nil
}

func (e *OAuth2APIEndpoint) URL() (*url.URL, error) {

	raw := fmt.Sprintf("%s://%s/%s", e.Scheme, e.Host, e.Path)

	url, err := url.Parse(raw)

	if err != nil {
		return nil, err
	}

	return url, err
}

func (e *OAuth2APIEndpoint) NewRequest(params *url.Values) (*http.Request, error) {

	url, err := e.URL()

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url.String(), nil)

	if err != nil {
		return nil, err
	}

	params.Set("access_token", e.Token)
	req.URL.RawQuery = (*params).Encode()

	return req, nil
}
