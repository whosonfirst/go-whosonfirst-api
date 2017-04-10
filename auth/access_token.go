package auth

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-api"
	"net/http"
)

type AccessTokenAuthentication struct {
	api.APIAuthentication
	access_token string
}

func NewAccessTokenAuthentication(token string) (*AccessTokenAuthentication, error) {

	au := AccessTokenAuthentication{
		access_token: token,
	}

	return &au, nil
}

func (au *AccessTokenAuthentication) AppendAuthentication(req *http.Request) error {

	params := req.URL.Query()
	params.Set("access_token", au.access_token)
	req.URL.RawQuery = (*params).Encode()

	return nil
}
