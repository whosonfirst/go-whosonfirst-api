package writer

import (
	"bufio"
	"encoding/json"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/util"
	"io"
	_ "log"
	"sync"
)

type GeoJSONLSWriter struct {
	api.APIResultWriter
	mu     *sync.Mutex
	writer *bufio.Writer
}

func NewGeoJSONLSWriter(w io.Writer) (*GeoJSONLSWriter, error) {

	mu := new(sync.Mutex)

	wr := GeoJSONLSWriter{
		mu:     mu,
		writer: bufio.NewWriter(w),
	}

	return &wr, nil
}

func (wr *GeoJSONLSWriter) WriteResult(r api.APIPlacesResult) (int, error) {

	geojson, err := util.APIResultToGeoJSON(r)

	if err != nil {
		return 0, err
	}

	var tmp interface{}

	err = json.Unmarshal(geojson, &tmp)

	if err != nil {
		return 0, err
	}

	body, err := json.Marshal(tmp)

	if err != nil {
		return 0, err
	}

	wr.mu.Lock()
	defer wr.mu.Unlock()

	n, err := wr.Write(body)

	if err != nil {
		return n, err
	}

	wr.Write([]byte("\n"))

	return n, nil
}

func (wr *GeoJSONLSWriter) Write(p []byte) (int, error) {
	i, err := wr.writer.Write(p)

	if err == nil {
		wr.writer.Flush()
	}

	return i, err
}

func (wr *GeoJSONLSWriter) Close() error {
	return nil
}
