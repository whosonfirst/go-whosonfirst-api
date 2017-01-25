package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/result"
	_ "log"
)

type JSONPagination struct {
	page     int
	pages    int
	per_page int
	total    int
}

func (p JSONPagination) Page() int {
	return p.page
}

func (p JSONPagination) Pages() int {
	return p.pages
}

func (p JSONPagination) PerPage() int {
	return p.per_page
}

func (p JSONPagination) Total() int {
	return p.total
}

func (p JSONPagination) String() string {
	return fmt.Sprintf("total %d page %d/%d (%d per page)", p.Total(), p.Page(), p.Pages(), p.PerPage())
}

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
	raw []byte
}

func (rsp JSONResponse) Raw() []byte {
	return rsp.raw
}

func (rsp JSONResponse) String() string {
	return string(rsp.raw)
}

func (rsp JSONResponse) Stat() string {

	r := rsp.get("stat")
	return r.String()
}

func (rsp JSONResponse) Ok() (bool, api.APIError) {

	if rsp.Stat() == "ok" {
		return true, nil
	}

	// TO DO: support this stuff
	// {"meta":{"version":1,"status_code":429},"results":{"error":{"type":"QpsExceededError","message":"Queries per second exceeded: Queries exceeded (1 allowed)."}}}

	code := rsp.get("error.code")
	msg := rsp.get("error.message")

	err := JSONError{
		code:    code.Int(),
		message: msg.String(),
	}

	return false, &err
}

func (rsp JSONResponse) Results() ([]api.APIResult, error) {

	results := make([]api.APIResult, 0)

	_results := rsp.get("results")

	// TO DO: signal failed NewJSONResult

	_results.ForEach(func(key, value gjson.Result) bool {

		_result, err := result.NewJSONResult(value)

		if err != nil {
			return false
		}

		results = append(results, _result)
		return true
	})

	return results, nil
}

func (rsp JSONResponse) get(path string) gjson.Result {
	return gjson.GetBytes(rsp.raw, path)
}

func (rsp JSONResponse) Pagination() (api.APIPagination, error) {

	page := rsp.get("page")

	if !page.Exists() {
		return nil, errors.New("Response is not paginated")
	}

	pages := rsp.get("pages")
	per_page := rsp.get("per_page")
	total := rsp.get("total")

	pg := JSONPagination{
		page:     int(page.Int()),
		pages:    int(pages.Int()),
		per_page: int(per_page.Int()),
		total:    int(total.Int()),
	}

	return &pg, nil
}

func ParseJSONResponse(raw []byte) (*JSONResponse, error) {

	var stub interface{}
	err := json.Unmarshal(raw, &stub)

	if err != nil {
		return nil, err
	}

	rsp := JSONResponse{
		raw: raw,
	}

	return &rsp, nil
}
