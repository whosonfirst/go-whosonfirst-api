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
     // please write me
}

type APIResponse interface {
     // please write me
}
