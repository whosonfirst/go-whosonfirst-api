package api

import (
       "fmt"
       "net/http"
       "net/url"
)

type APIEndpoint interface {
     URL() (url.URL, error)
     NewRequest(*url.Values) (*http.Request, error)
}

type MapzenAPIEndpoint struct {
     APIEndpoint
     Scheme string
     Host string
     Path string
     Key string
}

func NewMapzenAPIEndpoint (key string) (*MapzenAPIEndpoint, error) {

     cl := MapzenAPIEndpoint{
     	Scheme: "https",
	Host: "whosonfirst-mapzen.com",
	Path: "",
	Key: key,
     }

     return &cl, nil
}

func (cl *MapzenAPIEndpoint) URL() (url.URL, error) {

     raw := fmt.Sprintf("%s://%s/%s", cl.Scheme, cl.Host, cl.Path)
     
     url, err := url.Parse(raw)
     
     if err != nil {
     		return nil, err
		}

		return 	url, err
}

func (cl *MapzenAPIEndpoint) NewRequest(params *url.Values) (*http.Request, error) {

     url, err := cl.URL()

     if err != nil {
     	return nil, err
     }
     
	req, err := http.NewRequest("POST", url.String(), nil)

	if err != nil {
		return nil, req_err
	}

	params.Set("api_key", cl.Key)
	req.URL.RawQuery = (*params).Encode()

	return req, nil	
}
