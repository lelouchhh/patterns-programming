package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Circuit func(ctx context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	failureCount := 0
	lastAttempt := time.Now()
	var m sync.RWMutex
	return func(ctx context.Context) (string, error) {
		m.RLock()
		d := failureCount - int(failureThreshold)
		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("circuit breaker: too many attempts")
			}
		}
		m.RUnlock()
		resp, err := circuit(ctx)
		m.Lock()
		defer m.Unlock()
		if err != nil {
			failureCount++
			return resp, err
		}
		failureCount = 0
		return resp, nil
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	serviceCall := func(ctx context.Context) (string, error) {
		return "error", errors.New("some sort of error")
	}
	defer cancel()
	breaker := Breaker(serviceCall, 5)
	for i := 0; i < 10; i++ {
		fmt.Println(breaker(ctx))
	}
}
