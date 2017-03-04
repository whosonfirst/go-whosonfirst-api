package util

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-api"
	"io"
	"io/ioutil"
	_ "log"
	"net/http"
)

func HTTPResponseToBytes(http_rsp *http.Response) ([]byte, error) {

	var body io.Reader
	var err error

	switch http_rsp.Header.Get("Content-Encoding") {

	case "gzip":

		body, err = gzip.NewReader(http_rsp.Body)

		if err != nil {
			return nil, err
		}

	default:
		body = http_rsp.Body
	}

	http_body, io_err := ioutil.ReadAll(body)

	if io_err != nil {
		return nil, io_err
	}

	return http_body, nil
}

func APIResultToGeoJSON(api_rsp api.APIResult) ([]byte, error) {

	uri := api_rsp.URI()

	http_rsp, err := http.Get(uri)

	if err != nil {
		return nil, err
	}

	defer http_rsp.Body.Close()

	if http_rsp.StatusCode != 200 {
		msg := fmt.Sprintf("Failed to retrieve %s because %s", uri, http_rsp.Status)
		return nil, errors.New(msg)
	}

	body, err := ioutil.ReadAll(http_rsp.Body)

	if err != nil {
		return nil, err
	}

	var tmp interface{}

	err = json.Unmarshal(body, &tmp)

	if err != nil {
		return nil, err
	}

	return body, nil
}
