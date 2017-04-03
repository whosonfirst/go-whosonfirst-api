package writer

// If this looks familiar it's because it's a copied-and-slightly-modified
// version of the source for io.MultiWriter (20170125/thisisaaronland)

import (
	"github.com/whosonfirst/go-whosonfirst-api"
	"log"
	"sync"
)

type APIResultMultiWriterAsync struct {
	api.APIResultMultiWriter
	writers  []api.APIResultWriter
	wg       *sync.WaitGroup
	throttle chan bool
}

func (mw *APIResultMultiWriterAsync) Write(r api.APIPlacesResult) (int, error) {

	for _, w := range mw.writers {

		<-mw.throttle

		mw.wg.Add(1)

		go func(w api.APIResultWriter, r api.APIPlacesResult, throttle chan bool, wg *sync.WaitGroup) {

			defer func() {
				throttle <- true
				wg.Done()
			}()

			_, err := w.WriteResult(r)

			if err != nil {
				log.Println(err)
			}

		}(w, r, mw.throttle, mw.wg)

	}

	return 1, nil
}

func (mw *APIResultMultiWriterAsync) Close() {

	mw.wg.Wait()

	for _, wr := range mw.writers {
		wr.Close()
	}
}

func NewAPIResultMultiWriterAsync(writers ...api.APIResultWriter) *APIResultMultiWriterAsync {

	w := make([]api.APIResultWriter, len(writers))
	copy(w, writers)

	count := 50 // sudo make me a knob
	throttle := make(chan bool, count)

	for i := 0; i < count; i++ {
		throttle <- true
	}

	wg := new(sync.WaitGroup)

	mw := APIResultMultiWriterAsync{
		writers:  w,
		throttle: throttle,
		wg:       wg,
	}

	return &mw
}
