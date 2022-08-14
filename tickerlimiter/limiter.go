package tickerlimiter

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

type Limiter interface {
	Handle(url string) error
}

type limiter struct {
	ctx                    context.Context
	Ticker                 *time.Ticker
	maxRequestsCount       int
	currentRequestsCounter int
}

func (l *limiter) Handle(url string) error {
	for {
		select {
		case <-l.Ticker.C:
			res, err := http.Get(url)
			if err != nil {
				return err
			}
			_, err = ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			l.currentRequestsCounter++
			if l.maxRequestsCount == l.currentRequestsCounter {
				return nil
			}
			// work...
		case <-l.ctx.Done():
			return nil
		}
	}
}

func NewLimiter(ctx context.Context, maxRequestsCount int, tickerDuration time.Duration) Limiter {
	l := &limiter{
		ctx:              ctx,
		Ticker:           time.NewTicker(tickerDuration),
		maxRequestsCount: maxRequestsCount,
	}
	return l
}
