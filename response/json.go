package response

import (
	"encoding/json"
	"github.com/tidwall/gjson"
        "github.com/whosonfirst/go-whosonfirst-api"
)

type APIError struct {
	Code    int64
	Message string
}

type JSONResponse struct {
     api.APIResponse     		  
	Raw []byte
}

func (rsp JSONResponse) String() string {
     return string(rsp.Raw)
}

func (rsp JSONResponse) Get(path string) gjson.Result {
	return gjson.GetBytes(rsp.Raw, path)
}

func (rsp JSONResponse) Stat() string {

	r := rsp.Get("stat")
	return r.String()
}

func (rsp JSONResponse) Ok() (bool, *APIError) {

	if rsp.Stat() == "ok" {
		return true, nil
	}

	return false, rsp.Error()
}

func (rsp JSONResponse) Error() *APIError {

	code := rsp.Get("error.code")
	msg := rsp.Get("error.message")

	err := APIError{
		Code:    code.Int(),
		Message: msg.String(),
	}

	return &err
}

func ParseJSONResponse(raw []byte) (*JSONResponse, error) {

	var stub interface{}
	err := json.Unmarshal(raw, &stub)

	if err != nil {
		return nil, err
	}

	rsp := JSONResponse{
		Raw: raw,
	}

	return &rsp, nil
}
