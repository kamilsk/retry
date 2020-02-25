// Package retry provides the most advanced interruptible mechanism
// to perform actions repetitively until successful.
package retry

import (
	"fmt"

	"github.com/kamilsk/retry/v5/strategy"
)

// Action defines a callable function that package retry can handle.
type Action func() error

// How is an alias for batch of Strategies.
//
//  how := retry.How{
//  	strategy.Limit(3),
//  }
//
type How []func(strategy.Breaker, uint, error) bool

// Do takes an action and performs it, repetitively, until successful.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Do(
	breaker strategy.Breaker,
	action func() error,
	strategies ...func(strategy.Breaker, uint, error) bool,
) error {
	var err error
	for attempt, should := uint(0), true; should; attempt++ {
		for i, repeat := 0, len(strategies); should && i < repeat; i++ {
			should = should && strategies[i](breaker, attempt, err)
		}
		select {
		case <-breaker.Done():
			return breaker.Err()
		default:
			if should {
				err = action()
			}
		}
		should = should && err != nil
	}
	return err
}

// DoAsync takes an action and performs it, repetitively, until successful.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func DoAsync(
	breaker strategy.Breaker,
	action func() error,
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
