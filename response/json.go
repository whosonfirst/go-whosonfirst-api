package response

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-api"
)

type JSONError struct {
	api.APIError
	code    int64
	message string
}

func (e JSONError) Code() int64 {
	return e.code
}

func (e JSONError) Message() string {
	return e.message
}

func (e JSONError) String() string {
	return fmt.Sprintf("%d %s", e.Code(), e.Message())
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

func (rsp JSONResponse) Ok() (bool, api.APIError) {

	if rsp.Stat() == "ok" {
		return true, nil
	}

	code := rsp.Get("error.code")
	msg := rsp.Get("error.message")

	err := JSONError{
		code:    code.Int(),
		message: msg.String(),
	}

	return false, &err
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
