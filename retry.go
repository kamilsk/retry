// Package retry provides the most advanced functional mechanism
// to perform actions repetitively until successful until successful.
package retry

import (
	"context"
	"errors"
	"sync/atomic"
)

// Retry takes action and performs it, repetitively, until successful.
// When it is done it releases resources associated with the Breaker.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Retry(
	breaker BreakCloser,
	action func(attempt uint) error,
	strategies ...func(attempt uint, err error) bool,
) error {
	if breaker != nil {
		defer breaker.Close()
	}
	return retry(breaker, action, strategies...)
}

// Try takes action and performs it, repetitively, until successful.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Try(
	breaker Breaker,
	action func(attempt uint) error,
	strategies ...func(attempt uint, err error) bool,
) error {
	return retry(breaker, action, strategies...)
}

// TryContext takes action and performs it, repetitively, until successful.
// It uses the Context as a Breaker to prevent unnecessary action execution.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func TryContext(
	ctx context.Context,
	action func(ctx context.Context, attempt uint) error,
	strategies ...func(attempt uint, err error) bool,
) error {
	return retry(ctx, currying(ctx, action), strategies...)
}

// IsInterrupted checks that the error is related to the Breaker interruption
// on Retry call.
func IsInterrupted(err error) bool {
	return err == errInterrupted
}

// IsRecovered checks that the error is related to unhandled Action's panic
// and returns an original cause of panic.
func IsRecovered(err error) (interface{}, bool) {
	if h, is := err.(panicHandler); is {
		return h.recovered, true
	}
	return nil, false
}

var (
	errPanic       = errors.New("unhandled action's panic")
	errInterrupted = errors.New("operation was interrupted")
)

func currying(ctx context.Context, action func(context.Context, uint) error) func(uint) error {
	return func(attempt uint) error {
		return action(ctx, attempt)
	}
}

func retry(
	breaker Breaker,
	action func(attempt uint) error,
	strategies ...func(attempt uint, err error) bool,
) error {
	if (breaker == nil || breaker.Done() == nil) && len(strategies) == 0 {
		return action(0)
	}
	var (
		interrupted  uint32
		interruption <-chan struct{}
	)
	if breaker != nil {
		interruption = breaker.Done()
	}

	done := make(chan error, 1)
	go func() {
		var err error

		defer close(done)
		defer func() { done <- err }()
		defer panicHandler{}.recover(&err)

		for attempt := uint(0); shouldAttempt(attempt, err, strategies...) &&
			!atomic.CompareAndSwapUint32(&interrupted, 1, 0); attempt++ {

			err = action(attempt)
		}
	}()

	select {
	case <-interruption:
		atomic.CompareAndSwapUint32(&interrupted, 0, 1)
		return errInterrupted
	case err := <-done:
		return err
	}
}

// shouldAttempt evaluates the provided strategies with the given attempt to
// determine if the Retry loop should make another attempt.
func shouldAttempt(attempt uint, err error, strategies ...func(uint, error) bool) bool {
	should := attempt == 0 || err != nil

	for i, repeat := 0, len(strategies); should && i < repeat; i++ {
		should = should && strategies[i](attempt, err)
	}

	return should
}

type panicHandler struct {
	error
	recovered interface{}
}

func (panicHandler) recover(err *error) {
	if r := recover(); r != nil {
		*err = panicHandler{errPanic, r}
	}
}
