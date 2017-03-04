package response

import (
       "bytes"
	_ "errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/result"
	"github.com/whosonfirst/go-whosonfirst-api/util"	
	"github.com/whosonfirst/go-whosonfirst-csv"
	"io"
	_ "log"
	"net/http"
	"strconv"
)

type CSVResponse struct {
	api.APIResponse
	raw []byte
	pagination CSVPagination
}

type CSVPagination struct {
	page       int
	pages      int
	per_page   int
	total      int
	cursor     string
	next_query string
}

func (p CSVPagination) Page() int {
	return p.page
}

func (p CSVPagination) Pages() int {
	return p.pages
}

func (p CSVPagination) PerPage() int {
	return p.per_page
}

func (p CSVPagination) Total() int {
	return p.total
}

func (p CSVPagination) Cursor() string {
	return p.cursor
}

func (p CSVPagination) NextQuery() string {
	return p.next_query
}

func (p CSVPagination) String() string {
	return fmt.Sprintf("total %d page %d/%d (%d per page) cursor %s next_query %s", p.Total(), p.Page(), p.Pages(), p.PerPage(), p.Cursor(), p.NextQuery())
}

/*
func (rsp CSVResponse) ToBuffer() (*bytes.Buffer) {

     fieldnames := rsp.reader.Fieldname()
     buf := new(bytes.Buffer)

     writer := csv.NewDictWriter(buf, fieldname)
     writer.WriteHeader()

     for {
     	 row, err := rsp.reader.Read()

    	 if err == io.EOF {
            break
	 }

	 writer.WriteRow(row)
    }

    return buf
}
*/

func (rsp CSVResponse) Raw() []byte {
	return rsp.raw
}

func (rsp CSVResponse) String() string {
	return string(rsp.raw)
}

func (rsp CSVResponse) Ok() (bool, api.APIError) {
	return true, nil
}

func (rsp CSVResponse) Pagination() (api.APIPagination, error) {
	return rsp.pagination, nil
}

func (rsp CSVResponse) Results() ([]api.APIResult, error) {

	results := make([]api.APIResult, 0)
	
	byte_reader := bytes.NewReader(rsp.raw)
	reader, err := csv.NewDictReader(byte_reader)

	if err != nil {
	   return results, err
	}

	for {

		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return results, err
		}

		_result, err := result.NewCSVResult(row)

		if err != nil {
			return results, err
		}

		results = append(results, _result)
	}

	return results, nil
}

func ParseCSVResponse(http_rsp *http.Response) (*CSVResponse, error) {

	raw, err := util.HTTPResponseToBytes(http_rsp)

	if err != nil {
		return nil, err
	}

	header := http_rsp.Header

	str_page := header.Get("X-api-pagination-page")
	str_pages := header.Get("X-api-pagination-pages")
	str_per_page := header.Get("X-api-pagination-per-page")
	str_total := header.Get("X-api-pagination-total")
	cursor := header.Get("X-api-pagination-cursor")
	next_query := header.Get("X-api-pagination-next-query")

	page, err := strconv.Atoi(str_page)

	if err != nil {
		return nil, err
	}

	pages, err := strconv.Atoi(str_pages)

	if err != nil {
		return nil, err
	}

	per_page, err := strconv.Atoi(str_per_page)

	if err != nil {
		return nil, err
	}

	total, err := strconv.Atoi(str_total)

	if err != nil {
		return nil, err
	}

	pg := CSVPagination{
		page:       page,
		pages:      pages,
		per_page:   per_page,
		total:      total,
		cursor:     cursor,
		next_query: next_query,
	}

	rsp := CSVResponse{
		raw:        raw,
		pagination: pg,
	}

	return &rsp, nil
}
