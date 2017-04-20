package writer

import (
	"github.com/whosonfirst/go-whosonfirst-api"
	"io"
	_ "log"
	"sync"
)

type FileListWriter struct {
	api.APIResultWriter
	features int
	mu       *sync.Mutex
	writer   io.Writer
}

func NewFileListWriter(w io.Writer) (*FileListWriter, error) {

	mu := new(sync.Mutex)

	wr := FileListWriter{
		features: 0,
		mu:       mu,
		writer:   w,
	}

	return &wr, nil
}

func (wr *FileListWriter) WriteResult(r api.APIResult) (int, error) {

	wr.mu.Lock()
	defer wr.mu.Unlock()

	path := r.Path()

	n, err := wr.Write([]byte(path + "\n"))

	if err != nil {
		return n, err
	}

	wr.features += 1
	return n, nil
}

func (wr *FileListWriter) Write(p []byte) (int, error) {
	return wr.writer.Write(p)
}

func (wr *FileListWriter) Close() error {
	return nil
}
