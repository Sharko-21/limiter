package main

import (
	"context"
	"log"
	"time"

	"ratelimiter/tickerlimiter"
)

func main() {
	/*l := limiter.NewLimiter(limiter.Params{
		LimiterTime:        1 * time.Millisecond,
		MaxRequestsPerTime: 2,
	})
	for {
		if _, err := l.Get("https://google.com"); err != nil {
			log.Fatal(err)
		}
	}*/
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	l2 := tickerlimiter.NewLimiter(ctx, 5, 1*time.Second)
	if err := l2.Handle("https://google.com"); err != nil {
		log.Fatal(err)
	}
}
