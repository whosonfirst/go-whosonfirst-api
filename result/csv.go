package result

import (
	"bytes"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"sort"
	"strconv"
)

type CSVResult struct {
	api.APIPlacesResult
	result map[string]string
}

func NewCSVResult(result map[string]string) (*CSVResult, error) {

	r := CSVResult{
		result: result,
	}

	return &r, nil
}

func (r CSVResult) String(flags ...api.APIResultFlag) string {

	fieldnames := make([]string, 0)

	for k, _ := range r.result {
		fieldnames = append(fieldnames, k)
	}

	sort.Strings(fieldnames)

	buf := new(bytes.Buffer)

	writer, err := csv.NewDictWriter(buf, fieldnames)

	if err != nil {
		return ""
	}

	// this sucks... (2070304/thisisaaronland)
	// APIResultFlag is all wet paint... still trying to work it out

	if len(flags) > 0 && flags[0].Key() == "header" && flags[0].Bool() {
		writer.WriteHeader()
	}

	writer.WriteRow(r.result)
	return buf.String()
}

func (r CSVResult) WOFId() int64 {

	str_id, _ := r.result["wof_id"]
	id, _ := strconv.Atoi(str_id)

	return int64(id)
}

func (r CSVResult) WOFParentId() int64 {

	str_id, _ := r.result["wof_parent_id"]
	id, _ := strconv.Atoi(str_id)

	return int64(id)
}

func (r CSVResult) WOFName() string {
	name, _ := r.result["wof_name"]
	return name
}

func (r CSVResult) WOFPlacetype() string {
	placetype, _ := r.result["wof_placetype"]
	return placetype
}

func (r CSVResult) WOFCountry() string {
	country, _ := r.result["wof_country"]
	return country
}

func (r CSVResult) WOFRepo() string {
	repo, _ := r.result["wof_repo"]
	return repo
}

func (r CSVResult) Path() string {
	path, _ := uri.Id2RelPath(r.WOFId())
	return path
}

func (r CSVResult) URI() string {
	uri, _ := uri.Id2AbsPath("https://whosonfirst.mapzen.com/data", r.WOFId())
	return uri
}
