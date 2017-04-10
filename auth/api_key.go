package auth

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-api"
	"net/http"
)

type APIKeyAuthentication struct {
	api.APIAuthentication
	api_key string
}

func NewAPIKeyAuthentication(key string) (*APIKeyAuthentication, error) {

	au := APIKeyAuthentication{
		api_key: token,
	}

	return &au, nil
}

func (au *APIKeyAuthentication) AppendAuthentication(req *http.Request) error {

	params := req.URL.Query()
	params.Set("api_key", au.api_key)
	req.URL.RawQuery = (*params).Encode()

	return nil
}
