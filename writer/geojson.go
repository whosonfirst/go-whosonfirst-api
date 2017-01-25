package writer

import (
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/util"
	"os"
	"sync"
)

type GeoJSONWriter struct {
	features int
	mu       *sync.Mutex
}

func NewGeoJSONWriter() (*GeoJSONWriter, error) {

	mu := new(sync.Mutex)

	w := GeoJSONWriter{
		features: 0,
		mu:       mu,
	}

	w.Write([]byte(`{"type":"FeatureCollection", "features":[`))
	w.features = 0

	return &w, nil
}

func (w *GeoJSONWriter) WriteResult(r api.APIResult) (int, error) {

	json, err := util.APIResultToGeoJSON(r)

	if err != nil {
		return 0, err
	}

	return w.Write(json)
}

func (w *GeoJSONWriter) Write(p []byte) (int, error) {

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.features > 0 {
		os.Stdout.Write([]byte(`,`))
	}

	i, err := os.Stdout.Write(p)

	if err == nil {

		w.features += 1
	}

	return i, err
}

func (w *GeoJSONWriter) Close() error {
	w.features = 0
	w.Write([]byte(`]}`))
	return nil
}
