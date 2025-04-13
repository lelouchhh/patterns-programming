package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type effector func(ctx context.Context) (string, error)

func retry(effector effector, retries int, delay time.Duration) effector {
	return func(ctx context.Context) (string, error) {
		for n := 0; ; n++ {
			resp, err := effector(context.Background())
			if err == nil || n >= retries {
				return resp, err
			}
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}

func main() {
	ctx := context.Background()
	var i int
	fn := retry(func(ctx context.Context) (string, error) {
		if i < 4 {
			i++
			return "error!", errors.New("invalid i")
		}
		return "alright!", nil
	}, 10, time.Second)
	fmt.Println(fn(ctx))
}
