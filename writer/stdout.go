package writer

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-api"
	"os"
)

type StdoutWriter struct {
	api.APIResultWriter
}

func NewStdoutWriter() (*StdoutWriter, error) {

	wr := StdoutWriter{}
	return &wr, nil
}

func (wr *StdoutWriter) WriteResult(r api.APIPlacesResult) (int, error) {

	text := fmt.Sprintf("%d %s %s\n", r.WOFId(), r.WOFPlacetype(), r.WOFName())
	return wr.Write([]byte(text))
}

func (wr *StdoutWriter) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (wr *StdoutWriter) Close() error {
	return nil
}
