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

func (fl StdoutFlags) Filehandles() ([]io.Writer, error) {

	writers := []io.Writer{os.Stdout}
	return writers, nil
}
