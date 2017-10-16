package flags

import (
	"io"
	"os"
	"strings"
)

type FileHandleFlags struct {
	ResultWriterFlags
	flags []string
}

func (fl *FileHandleFlags) String() string {
	return strings.Join(fl.flags, "\n")
}

func (fl *FileHandleFlags) Set(value string) error {
	fl.flags = append(fl.flags, value)
	return nil
}

func (fl FileHandleFlags) FileHandles() ([]io.Writer, error) {

	writers := make([]io.Writer, 0)

	for _, path := range fl.flags {

		if path == "-" {
			fh := os.Stdout
			writers = append(writers, fh)
		} else {

			fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)

			if err != nil {
				return nil, err
			}

			writers = append(writers, fh)
		}
	}

	return writers, nil
}
