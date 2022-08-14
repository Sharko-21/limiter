package limiter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Limiter interface {
	Get(url string) ([]byte, error)
}

type Params struct {
	MaxRequestsPerTime int
	LimiterTime        time.Duration
	LastRequestAt      time.Time

	requestsCounter int
	controller      chan struct{}
	starter         chan struct{}
}
type limiter struct {
	params Params
}

func (l *limiter) Get(url string) ([]byte, error) {
	l.params.controller <- struct{}{}
	<-l.params.starter
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("Запрос прошел")
	return responseData, err
}

func (l *limiter) startCommunicator() {
	for range l.params.controller {
		diff := time.Since(l.params.LastRequestAt)
		if l.params.requestsCounter == l.params.MaxRequestsPerTime && diff < l.params.LimiterTime {
			time.Sleep(l.params.LimiterTime - diff)
			l.params.requestsCounter = 0
		}
		l.params.LastRequestAt = time.Now()
		l.params.requestsCounter++
		l.params.starter <- struct{}{}
	}
}

func NewLimiter(p Params) Limiter {
	p.controller = make(chan struct{})
	p.starter = make(chan struct{})
	l := &limiter{
		params: p,
	}
	go l.startCommunicator()
	return l
}
