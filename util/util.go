package util

import (
	"encoding/json"
	"errors"
	"fmt"
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
