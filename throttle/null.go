package throttle

import (
	"github.com/whosonfirst/go-whosonfirst-api"
)

type NullThrottle struct {
	api.APIThrottle
}

func NewNullThrottle() (api.APIThrottle, error) {

	thr := NullThrottle{}

	return &thr, nil
}

func (thr *NullThrottle) RateLimit() chan bool {

	ch := make(chan bool, 1)

	go func() {
		ch <- true
		return
	}()

	return ch
}
