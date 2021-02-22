// Package retry provides the most advanced interruptible mechanism
// to perform actions repetitively until successful.
// The retry based on https://github.com/Rican7/retry but fully reworked
// and focused on integration with the https://github.com/kamilsk/breaker
// and the built-in https://pkg.go.dev/context package.
package retry

import (
	"context"
	"fmt"

	"github.com/kamilsk/retry/v5/strategy"
)

// Action defines a callable function that package retry can handle.
type Action = func(context.Context) error

// Error defines a string-based error without a different root cause.
type Error string

// Error returns a string representation of an error.
func (err Error) Error() string { return string(err) }

// Unwrap always returns nil means that an error doesn't have other root cause.
func (err Error) Unwrap() error { return nil }

// How is an alias for batch of Strategies.
//
//  how := retry.How{
//  	strategy.Limit(3),
//  }
//
type How = []func(strategy.Breaker, uint, error) bool

// Do takes the action and performs it, repetitively, until successful.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Do(
	breaker strategy.Breaker,
	action func(context.Context) error,
	strategies ...func(strategy.Breaker, uint, error) bool,
) error {
	var (
		err   error = Error("have no any try")
		clean error
	)

	ctx, is := breaker.(context.Context)
	if !is {
		ctx = lite{context.Background(), breaker.Done()}
	}

	for attempt, should := uint(0), true; should; attempt++ {
		clean = unwrap(err)
		for i, repeat := 0, len(strategies); should && i < repeat; i++ {
			should = should && strategies[i](breaker, attempt, clean)
		}

		select {
		case <-breaker.Done():
			return breaker.Err()
		default:
			if should {
				err = action(ctx)
			}
		}

		should = should && err != nil
	}

	return err
}

// Go takes the action and performs it, repetitively, until successful.
// It differs from the Do method in that it performs the action in a goroutine.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Go(
	breaker strategy.Breaker,
	action func(context.Context) error,
	strategies ...func(strategy.Breaker, uint, error) bool,
) error {
	done := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("retry: unexpected panic: %#v", r)
				}
				done <- err
			}
			close(done)
		}()
		done <- Do(breaker, action, strategies...)
	}()

	select {
	case <-breaker.Done():
		return breaker.Err()
	case err := <-done:
		return err
	}
}
