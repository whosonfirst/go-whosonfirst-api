package api

// PLEASE MAKE ME SUPPORT (not-just-json)

import (
       "encoding/json"
       "github.com/tidwall/gjson"
)

type APIError struct {
	Code    int64
	Message string
}

type APIResponse struct {
	Raw    []byte
}

func (rsp APIResponse) Get(path string) gjson.Result {
     return gjson.GetBytes(rsp.Raw, path)
}

func (rsp APIResponse) Stat() string {

     	r := rsp.Get("stat")
	return r.String()
}

func (rsp APIResponse) Ok() (bool, *APIError) {

     if rsp.Stat() == "ok" {
		return true, nil
	}

	return false, rsp.Error()
}

func (rsp APIResponse) Error() *APIError {

	code, _ = rsp.Get("error.code")
	msg, _ = rsp.Get("error.message")

	err := APIError{
	    Code: code.Int64(),
	    Message: msg.String(),
	}
	
	return &err
}

func ParseAPIResponse(raw []byte) (*APIResponse, error) {

     			  var stub interface{}
        err := json.Unmarshal(raw, &stub)

	if err != nil {
		return nil, err
	}
	
	rsp := APIResponse{
		Raw:    raw,
	}

	return &rsp, nil
}
