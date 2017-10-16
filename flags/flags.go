package flags

import (
	"github.com/whosonfirst/go-whosonfirst-api"
	"io"
)

type ResultWriterFlags interface {
	FileHandles() ([]io.Writer, error)
}

type ResultWriterFunc func(io.Writer) (api.APIResultWriter, error)

func ResultWriterFlagsToResultWriters(flags_map map[ResultWriterFlags]ResultWriterFunc) ([]api.APIResultWriter, error) {

	writers := make([]api.APIResultWriter, 0)

	for writers_flags, writers_func := range flags_map {

		filehandles, err := writers_flags.FileHandles()

		if err != nil {
			return nil, err
		}

		for _, fh := range filehandles {

			wr, err := writers_func(fh)

			if err != nil {
				return nil, err
			}

			writers = append(writers, wr)
		}
	}

	return writers, nil
}
