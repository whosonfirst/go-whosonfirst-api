package result

import (
	"bytes"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"strconv"
)

type MetaResult struct {
	api.APIPlacesResult
	result map[string]string
}

func NewMetaResult(result map[string]string) (*MetaResult, error) {

	r := MetaResult{
		result: result,
	}

	return &r, nil
}

func (r MetaResult) String(flags ...api.APIResultFlag) string {

	fieldnames := make([]string, 0)

	for k, _ := range r.result {
		fieldnames = append(fieldnames, k)
	}

	buf := new(bytes.Buffer)

	writer, err := csv.NewDictWriter(buf, fieldnames)

	if err != nil {
		return ""
	}

	writer.WriteHeader()
	writer.WriteRow(r.result)

	return buf.String()
}

func (r MetaResult) WOFId() int64 {

	str_id, _ := r.result["id"]
	id, _ := strconv.Atoi(str_id)

	return int64(id)
}

func (r MetaResult) WOFParentId() int64 {

	str_id, _ := r.result["parent_id"]
	id, _ := strconv.Atoi(str_id)

	return int64(id)
}

func (r MetaResult) WOFName() string {
	name, _ := r.result["name"]
	return name
}

func (r MetaResult) WOFPlacetype() string {
	placetype, _ := r.result["placetype"]
	return placetype
}

func (r MetaResult) WOFCountry() string {
	country, _ := r.result["wof_country"]
	return country
}

func (r MetaResult) WOFRepo() string {
	return ""
}

func (r MetaResult) Path() string {
	path, _ := uri.Id2RelPath(r.WOFId())
	return path
}

func (r MetaResult) URI() string {
	uri, _ := uri.Id2AbsPath("https://whosonfirst.mapzen.com/data", r.WOFId())
	return uri
}
