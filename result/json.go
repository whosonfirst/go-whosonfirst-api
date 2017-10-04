package result

import (
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type JSONResult struct {
	api.APIPlacesResult
	result gjson.Result
}

func NewJSONResult(result gjson.Result) (*JSONResult, error) {

	r := JSONResult{
		result: result,
	}

	return &r, nil
}

func (r JSONResult) String(flags ...api.APIResultFlag) string {
	return r.result.String()
}

func (r JSONResult) WOFId() int64 {
	return r.get("wof:id").Int()
}

func (r JSONResult) WOFParentId() int64 {
	return r.get("wof:parent_id").Int()
}

func (r JSONResult) WOFName() string {
	return r.get("wof:name").String()
}

func (r JSONResult) WOFPlacetype() string {
	return r.get("wof:placetype").String()
}

func (r JSONResult) WOFCountry() string {
	return r.get("wof:country").String()
}

func (r JSONResult) WOFRepo() string {
	return r.get("wof:repo").String()
}

func (r JSONResult) Path() string {
	path, _ := uri.Id2RelPath(r.WOFId())
	return path
}

func (r JSONResult) URI() string {
	uri, _ := uri.Id2AbsPath("https://whosonfirst.mapzen.com/data", r.WOFId())
	return uri
}

func (r JSONResult) get(path string) gjson.Result {
	return r.result.Get(path)
}
