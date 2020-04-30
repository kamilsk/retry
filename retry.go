// Package retry provides the most advanced interruptible mechanism
// to perform actions repetitively until successful.
package retry

import (
	"context"
	"fmt"

	"github.com/kamilsk/retry/v5/strategy"
)

// Action defines a callable function that package retry can handle.
type Action func(context.Context) error

// How is an alias for batch of Strategies.
//
//  how := retry.How{
//  	strategy.Limit(3),
//  }
//
type How []func(strategy.Breaker, uint, error) bool

// Do takes the action and performs it, repetitively, until successful.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Do(
	breaker strategy.Breaker,
	action func(context.Context) error,
	strategies ...func(strategy.Breaker, uint, error) bool,
) error {
	var err, clean error
	ctx, cancel := context.WithCancel(context.Background())
	for attempt, should := uint(0), true; should; attempt++ {
		clean = unwrap(err)
		for i, repeat := 0, len(strategies); should && i < repeat; i++ {
			should = should && strategies[i](breaker, attempt, clean)
		}
		select {
		case <-breaker.Done():
			cancel()
			return breaker.Err()
		default:
			if should {
				err = action(ctx)
			}
		}
		should = should && err != nil
	}
	cancel()
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

// equal to go.octolab.org/errors.Unwrap
func unwrap(err error) error {
	// compatible with github.com/pkg/errors
	type causer interface {
		Cause() error
	}
	// compatible with built-in errors since 1.13
	type wrapper interface {
		Unwrap() error
	}

	for err != nil {
		layer, is := err.(wrapper)
		if is {
			err = layer.Unwrap()
			continue
		}
		cause, is := err.(causer)
		if is {
			err = cause.Cause()
			continue
		}
		break
	}
	return err
}
