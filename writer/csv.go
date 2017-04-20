package writer

import (
	"github.com/whosonfirst/go-whosonfirst-api"
	"io"
	_ "log"
	"sync"
)

type CSVWriter struct {
	api.APIResultWriter
	features int
	mu       *sync.Mutex
	writer   io.Writer
}

func NewCSVWriter(w io.Writer) (*CSVWriter, error) {

	mu := new(sync.Mutex)

	wr := CSVWriter{
		features: 0,
		mu:       mu,
		writer:   w,
	}

	return &wr, nil
}

func (wr *CSVWriter) WriteResult(r api.APIPlacesResult) (int, error) {

	wr.mu.Lock()
	defer wr.mu.Unlock()

	var body string

	if wr.features == 0 {
		fl := api.NewAPIResultBooleanFlag("header", true)
		body = r.String(fl)
	} else {
		body = r.String()
	}

	n, err := wr.Write([]byte(body))

	if err != nil {
		return n, err
	}

	wr.features += 1
	return n, nil
}

func (wr *CSVWriter) Write(p []byte) (int, error) {
	return wr.writer.Write(p)
}

func (wr *CSVWriter) Close() error {
	return nil
}
