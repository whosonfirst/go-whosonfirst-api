package writer

// If this looks familiar it's because it's a copied-and-slightly-modified
// version of the source for io.MultiWriter (20170125/thisisaaronland)

import (
	"github.com/whosonfirst/go-whosonfirst-api"
	"log"
	"sync"
)

type APIResultAsyncWriter struct {
	api.APIResultFooWriter
	writers  []api.APIResultWriter
	wg       *sync.WaitGroup
	throttle chan bool
}

func (t *APIResultAsyncWriter) Write(r api.APIResult) (int, error) {

	for _, w := range t.writers {

		<-t.throttle

		t.wg.Add(1)

		go func(w api.APIResultWriter, r api.APIResult, throttle chan bool, wg *sync.WaitGroup) {

			defer func() {
				throttle <- true
				wg.Done()
			}()

			_, err := w.WriteResult(r)

			if err != nil {
				log.Println(err)
			}

		}(w, r, t.throttle, t.wg)

	}

	return 1, nil
}

func (t *APIResultAsyncWriter) Close() {

	t.wg.Wait()

	for _, wr := range t.writers {
		wr.Close()
	}
}

func NewAPIResultAsyncWriter(writers ...api.APIResultWriter) *APIResultAsyncWriter {

	w := make([]api.APIResultWriter, len(writers))
	copy(w, writers)

	count := 50 // sudo make me a knob
	throttle := make(chan bool, count)

	for i := 0; i < count; i++ {
		throttle <- true
	}

	wg := new(sync.WaitGroup)

	async := APIResultAsyncWriter{
		writers:  w,
		throttle: throttle,
		wg:       wg,
	}

	return &async
}
