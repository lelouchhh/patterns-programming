package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Circuit func(ctx context.Context) (string, error)

func Debounce(circuit Circuit, d time.Duration) Circuit {
	var threshold time.Time
	var result string
	var err error
	var m sync.Mutex
	return func(ctx context.Context) (string, error) {
		m.Lock()
		defer func() {
			threshold = time.Now().Add(d)
			m.Unlock()
		}()
		if time.Now().Before(threshold) {
			return result, err
		}
		return circuit(ctx)
	}
}

func main() {
	res := Debounce(func(ctx context.Context) (string, error) {
		return "fetching data...", nil
	}, time.Duration(10*time.Microsecond))
	for i := 0; i <= 10; i++ {
		fmt.Println(res(context.TODO()))
	}

}
