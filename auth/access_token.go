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

func New AccessTokenAuthentication(token string) (*AccessTokenAuthentication, error) {

     au := AccessTokenAuthentication{
     	access_token: token,
     }

     return &au, nil
}

func (*au AccessTokenAuthentication) SetAuthentication(token string) error {
     au.access_token = token
     return nil
}

func (*au AccessTokenAuthentication) AddAuthentication(req *http.Request) error {

     return errors.New("Please write me")
}
