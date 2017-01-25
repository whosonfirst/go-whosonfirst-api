package writer

import (
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/util"
	"io"
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

	gj := GeoJSONWriter{
		features: 0,
		mu:       mu,
		writer:   w,
	}

	gj.Write([]byte(`{"type":"FeatureCollection", "features":[`))
	gj.features = 0

	return &gj, nil
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
		w.writer.Write([]byte(`,`))
	}

	i, err := w.writer.Write(p)

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
