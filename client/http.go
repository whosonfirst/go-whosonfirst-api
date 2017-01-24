package client

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/response"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type HTTPClient struct {
	api.APIClient
	endpoint api.APIEndpoint
}

func NewHTTPClient(endpoint api.APIEndpoint) (*HTTPClient, error) {

	cl := HTTPClient{
		endpoint: endpoint,
	}

	return &cl, nil
}

func (client *HTTPClient) DefaultArgs() *url.Values {
	args := url.Values{}
	return &args
}

func (client *HTTPClient) ExecuteMethod(method string, params *url.Values) (api.APIResponse, error) {

	params.Set("method", method)

	http_req, err := client.endpoint.NewRequest(params)

	if err != nil {
		return nil, err
	}

	http_req.Header.Add("Accept-Encoding", "gzip")

	http_client := &http.Client{}
	http_rsp, http_err := http_client.Do(http_req)

	if http_err != nil {
		return nil, http_err
	}

	defer http_rsp.Body.Close()

	http_body, io_err := ioutil.ReadAll(http_rsp.Body)

	if io_err != nil {
		return nil, io_err
	}

	// to do: support other formats...

	rsp, parse_err := response.ParseJSONResponse(http_body)

	if parse_err != nil {
		return nil, parse_err
	}

	return rsp, nil
}

func (client *HTTPClient) ExecuteMethodPaginated(method string, params *url.Values) error {

	pages := 0
	page := 1

	for {

		params.Set("page", string(page))

		log.Println("GO", page, pages)

		rsp, err := client.ExecuteMethod(method, params)

		if err != nil {
			return err
		}

		_, api_err := rsp.Ok()

		if api_err != nil {
			return errors.New(api_err.String())
		}

		if pages == 0 {

			pg, err := rsp.Pagination()

			log.Println(pg)

			if err != nil {
				return err
			}

			pages = pg.Pages()
		}

		/*
			cb_err := callback(rsp)

			if cb_err != nil {
				return cb_err
			}
		*/

		if page >= pages {
			log.Println("BREAK", page, pages)
			break
		}

		page += 1
	}

	return nil
}
