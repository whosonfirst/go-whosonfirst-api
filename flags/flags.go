package flags

import (
	"github.com/whosonfirst/go-whosonfirst-api"
	"io"
)

type ResultWriterFlags interface {
	FileHandles() ([]io.Writer, error)
}

type ResultWriterFunc func(io.Writer) (api.APIResultWriter, error)
