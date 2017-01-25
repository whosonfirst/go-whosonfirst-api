package writer

// If this looks familiar it's because it's a copied-and-slightly-modified
// version of the source for io.MultiWriter (20170125/thisisaaronland)

import (
	"github.com/whosonfirst/go-whosonfirst-api"
)

type APIResultMultiWriter struct {
	writers []api.APIResultWriter
}

func (t *APIResultMultiWriter) Write(r api.APIResult) (int, error) {

	var p int

	for _, w := range t.writers {

		n, err := w.WriteResult(r)

		if err != nil {
			return p, err
		}

		p += n
	}

	return p, nil
}

func NewAPIResultMultiWriter(writers ...api.APIResultWriter) *APIResultMultiWriter {
	w := make([]api.APIResultWriter, len(writers))
	copy(w, writers)
	return &APIResultMultiWriter{w}
}
