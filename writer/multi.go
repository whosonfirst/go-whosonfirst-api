package writer

// If this looks familiar it's because it's a copied-and-slightly-modified
// version of the source for io.MultiWriter (20170125/thisisaaronland)

import (
	"github.com/whosonfirst/go-whosonfirst-api"
	_ "log"
)

type APIResultMultiWriterSync struct {
	api.APIResultMultiWriter
	writers []api.APIResultWriter
}

func (mw *APIResultMultiWriterSync) Write(r api.APIPlacesResult) (int, error) {

	var p int

	for _, w := range mw.writers {

		n, err := w.WriteResult(r)

		if err != nil {
			return p, err
		}

		p += n
	}

	return p, nil
}

func (mw *APIResultMultiWriterSync) Close() {

	for _, wr := range mw.writers {
		wr.Close()
	}
}

func NewAPIResultMultiWriterSync(writers ...api.APIResultWriter) *APIResultMultiWriterSync {
	w := make([]api.APIResultWriter, len(writers))
	copy(w, writers)

	mw := APIResultMultiWriterSync{
		writers: w,
	}

	return &mw
}
