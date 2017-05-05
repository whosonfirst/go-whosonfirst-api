package writer

import (
	"encoding/json"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/util"
	"io"
	_ "log"
	"sync"
)

type GeoJSONWriter struct {
	api.APIResultWriter
	features int
	mu       *sync.Mutex
	writer   io.Writer
}

func NewGeoJSONWriter(w io.Writer) (*GeoJSONWriter, error) {

	mu := new(sync.Mutex)

	wr := GeoJSONWriter{
		features: 0,
		mu:       mu,
		writer:   w,
	}

	wr.Write([]byte(`{"type":"FeatureCollection", "features":[`))

	return &wr, nil
}

func (wr *GeoJSONWriter) WriteResult(r api.APIPlacesResult) (int, error) {

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

	if wr.features > 0 {
		wr.Write([]byte(`,`))
	}

	n, err := wr.Write(body)

	if err != nil {
		return n, err
	}

	wr.features += 1
	return n, nil
}

func (wr *GeoJSONWriter) Write(p []byte) (int, error) {
	return wr.writer.Write(p)
}

func (wr *GeoJSONWriter) Close() error {
	_, err := wr.Write([]byte(`]}`))
	return err
}
