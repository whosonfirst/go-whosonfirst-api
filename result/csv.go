package result

import (
	"bytes"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"strconv"
)

type CSVResult struct {
	api.APIResult
	result map[string]string
}

func NewCSVResult(result map[string]string) (*CSVResult, error) {

	r := CSVResult{
		result: result,
	}

	return &r, nil
}

func (r CSVResult) String() string {

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

func (r CSVResult) WOFId() int64 {

	str_id, _ := r.result["wof:id"]
	id, _ := strconv.Atoi(str_id)

	return int64(id)
}

func (r CSVResult) WOFParentId() int64 {

	str_id, _ := r.result["wof:parent_id"]
	id, _ := strconv.Atoi(str_id)

	return int64(id)
}

func (r CSVResult) WOFName() string {
	name, _ := r.result["wof:name"]
	return name
}

func (r CSVResult) WOFPlacetype() string {
	placetype, _ := r.result["wof:placetype"]
	return placetype
}

func (r CSVResult) WOFCountry() string {
	country, _ := r.result["wof:country"]
	return country
}

func (r CSVResult) WOFRepo() string {
	repo, _ := r.result["wof:repo"]
	return repo
}

func (r CSVResult) Path() string {
	path, _ := uri.Id2RelPath(int(r.WOFId()))
	return path
}

func (r CSVResult) URI() string {
	uri, _ := uri.Id2AbsPath("https://whosonfirst.mapzen.com/data", int(r.WOFId()))
	return uri
}
