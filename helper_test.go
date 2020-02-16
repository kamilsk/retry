package retry

import (
	"context"
	"time"
)

func delay(duration time.Duration) func(uint, error) bool {
	return func(attempt uint, _ error) bool {
		if 0 == attempt {
			time.Sleep(duration)
		}

		return true
	}
}

func limit(attemptLimit uint) func(uint, error) bool {
	return func(attempt uint, _ error) bool {
		return attempt < attemptLimit
	}
}

func newBreaker() *contextBreaker {
	ctx, cancel := context.WithCancel(context.Background())
	return &contextBreaker{ctx, cancel}
}

func newClosedBreaker() *contextBreaker {
	breaker := newBreaker()
	breaker.Close()
	return breaker
}

type contextBreaker struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (breaker *contextBreaker) Done() <-chan struct{} {
	return breaker.ctx.Done()
}

func (breaker *contextBreaker) Close() {
	breaker.cancel()
}

func (breaker *contextBreaker) Err() error {
	return breaker.ctx.Err()
}
