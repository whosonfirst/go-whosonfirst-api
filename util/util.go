package util

import (
	"encoding/json"
	"errors"
	"github.com/whosonfirst/go-whosonfirst-api"
	"io/ioutil"
	_ "log"
	"net/http"
)

func APIResultToGeoJSON(api_rsp api.APIResult) ([]byte, error) {

	uri := api_rsp.URI()

	http_rsp, err := http.Get(uri)

	if err != nil {
		return nil, err
	}

	defer http_rsp.Body.Close()

	if http_rsp.StatusCode != 200 {
		return nil, errors.New(http_rsp.Status)
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
