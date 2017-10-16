package throttle

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-api"
	"sync/atomic"
	"time"
)

type DefaultThrottle struct {
	api.APIThrottle
	QueriesPerSecond int32
	QueriesPerMinute int32
	QueriesPerHour   int32
	qpscount         int32
	qpmcount         int32
	qphcount         int32
}

func NewDefaultThrottle(ctx context.Context) (api.APIThrottle, error) {

	thr := DefaultThrottle{
		QueriesPerSecond: 6,
		QueriesPerMinute: 30,
		QueriesPerHour:   1000,
		qpscount:         0,
		qpmcount:         0,
		qphcount:         0,
	}

	ts := time.Tick(time.Second * 1)
	tm := time.Tick(time.Minute * 1)
	th := time.Tick(time.Hour * 1)

	watch := func(t <- chan time.Time, i int32) {

		for range t {

			atomic.StoreInt32(&i, 0)

			select {
			case <-ctx.Done():
				return
			}
		}
	}

	go watch(ts, thr.qpscount)
	go watch(tm, thr.qpmcount)
	go watch(th, thr.qphcount)

	return &thr, nil
}

func (thr *DefaultThrottle) RateLimit() chan bool {

	ch := make(chan bool, 1)

	go func() {

		for {

			qps := atomic.LoadInt32(&thr.qpscount)
			qpm := atomic.LoadInt32(&thr.qpmcount)
			qph := atomic.LoadInt32(&thr.qphcount)

			if qps >= thr.qpscount {
				time.Sleep(100 * time.Millisecond)
			} else if qpm >= thr.qpmcount {
				time.Sleep(200 * time.Millisecond)
			} else if qph >= thr.qphcount {
				time.Sleep(500 * time.Millisecond)
			} else {

				i := int32(1)
				atomic.AddInt32(&thr.qpscount, i)
				atomic.AddInt32(&thr.qpmcount, i)
				atomic.AddInt32(&thr.qphcount, i)

				ch <- true
				return
			}
		}
	}()

	return ch
}
