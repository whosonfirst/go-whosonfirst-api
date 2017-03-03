package writer

// If this looks familiar it's because it's a copied-and-slightly-modified
// version of the source for io.MultiWriter (20170125/thisisaaronland)

import (
	"github.com/whosonfirst/go-whosonfirst-api"
	"log"
	"sync"
)

type APIResultAsyncWriter struct {
	writers []api.APIResultWriter
	locks   []*sync.Mutex
	wg      *sync.WaitGroup
}

func (t *APIResultAsyncWriter) Write(r api.APIResult) (int, error) {

	for i, w := range t.writers {

		t.wg.Add(1)

		go func(w api.APIResultWriter, r api.APIResult, mu *sync.Mutex, wg *sync.WaitGroup) {

			defer wg.Done()

			mu.Lock()

			_, err := w.WriteResult(r)

			mu.Unlock()

			if err != nil {
				log.Println(err)
			}

		}(w, r, t.locks[i], t.wg)

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

	locks := make([]*sync.Mutex, 0)

	for range w {
		locks = append(locks, new(sync.Mutex))
	}

	wg := new(sync.WaitGroup)

	async := APIResultAsyncWriter{
		writers: w,
		locks:   locks,
		wg:      wg,
	}

	return &async
}
