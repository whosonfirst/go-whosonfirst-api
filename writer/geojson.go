package writer

import (
	"bufio"
	"bytes"
	"github.com/whosonfirst/go-whosonfirst-api"
	"github.com/whosonfirst/go-whosonfirst-api/util"
	"io"
	_ "log"
	"strings"
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

	json, err := util.APIResultToGeoJSON(r)

	if err != nil {
		return 0, err
	}

	wr.mu.Lock()
	defer wr.mu.Unlock()

	if wr.features > 0 {
		wr.Write([]byte(`,`))
	}

	n, err := wr.Write(json)

	if err != nil {
		return n, err
	}

	trim := false

	if trim {

		// sudo move me in to a separate package or something
		// (20170125/thisisaaronland)

		buf := bytes.NewBuffer(json)
		scanner := bufio.NewScanner(buf)

		n := 0

		for scanner.Scan() {

			str := scanner.Text()
			str = strings.Trim(str, "\r\n")
			str = strings.Trim(str, " ")

			i, err := wr.Write([]byte(str))

			if err != nil {
				return n, err
			}

			n += i
		}

		// end of sudo move me	in to a	separate package
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
