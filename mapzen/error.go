package mapzen

import (
	"encoding/json"
	"github.com/whosonfirst/go-whosonfirst-api/util"
	"net/http"
)

type Meta struct {
	Version    int `json:"version"`
	StatusCode int `json:"status_code"`
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Results struct {
	Error Error `json:"error"`
}

type Response struct {
	Meta    Meta    `json:"meta"`
	Results Results `json:"results"`
}

// {"meta":{"version":1,"status_code":403},"results":{"error":{"type":"KeyError","message":"'mapzen-xxxxpoo' is not a valid key."}}}

func ParseMapzenResponse(http_rsp *http.Response) (*Response, error) {

	raw, err := util.HTTPResponseToBytes(http_rsp)

	if err != nil {
		return nil, err
	}

	var rsp Response
	err = json.Unmarshal(raw, &rsp)

	if err != nil {
		return nil, err
	}

	return &rsp, nil
}
