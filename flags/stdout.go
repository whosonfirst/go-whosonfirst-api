package flags

import (
	"io"
	"os"
)

type StdoutFlags struct {
	ResultWriterFlags
	flag bool
}

func (fl *StdoutFlags) IsBoolFlag() bool {
	return true
}

func (fl *StdoutFlags) String() string {
	return "true"
}

func (fl *StdoutFlags) Set(value string) error {
	fl.flag = true
	return nil
}

func (fl StdoutFlags) FileHandles() ([]io.Writer, error) {

	writers := make([]io.Writer, 0)

	if fl.flag {
		fh := os.Stdout
		writers = append(writers, fh)
	}

	return writers, nil
}
