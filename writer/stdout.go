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

func (wr *StdoutWriter) WriteResult(r api.APIResult) (int, error) {

	text := fmt.Sprintf("%s is a %s with Who's On First ID %d\n", r.WOFName(), r.WOFPlacetype(), r.WOFId())
	return wr.Write([]byte(text))
}

func (wr *StdoutWriter) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (wr *StdoutWriter) Close() error {
	os.Stdout.Write([]byte("CLOSING TIME"))
	return nil
}
