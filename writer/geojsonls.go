package writer

import (
       "encoding/json"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/util"
	"io"
	_ "log"
	"sync"
)

type GeoJSONLSWriter struct {
	api.APIResultWriter
	mu       *sync.Mutex
	writer   io.Writer
}

func NewGeoJSONLSWriter(w io.Writer) (*GeoJSONLSWriter, error) {

	mu := new(sync.Mutex)

	wr := GeoJSONLSWriter{
		mu:       mu,
		writer:   w,
	}

	return &wr, nil
}

func (wr *GeoJSONLSWriter) WriteResult(r api.APIPlacesResult) (int, error) {

	geojson, err := util.APIResultToGeoJSON(r)

	if err != nil {
		return 0, err
	}

	body, err := json.Marshal(geojson)

	if err != nil {
		return 0, err
	}
	
	wr.mu.Lock()
	defer wr.mu.Unlock()
	
	n, err := wr.Write(body)

	if err != nil {
		return n, err
	}

return n, nil
}

func (wr *GeoJSONLSWriter) Write(p []byte) (int, error) {
	return wr.writer.Write(p)
}

func (wr *GeoJSONLSWriter) Close() error {
     return nil
}
