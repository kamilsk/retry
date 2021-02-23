package retry

import (
	"context"
	"fmt"
)

// Action defines a callable function that package retry can handle.
type Action = func(context.Context) error

// A Breaker carries a cancellation signal to interrupt an action execution.
//
// It is a subset of the built-in context and github.com/kamilsk/breaker interfaces.
type Breaker = interface {
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}
	// If Done is not yet closed, Err returns nil.
	// If Done is closed, Err returns a non-nil error.
	// After Err returns a non-nil error, successive calls to Err return the same error.
	Err() error
}

// How is an alias for batch of Strategies.
//
//  how := retry.How{
//  	strategy.Limit(3),
//  }
//
type How = []func(Breaker, uint, error) bool

// Do takes the action and performs it, repetitively, until successful.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Do(
	breaker Breaker,
	action func(context.Context) error,
	strategies ...func(Breaker, uint, error) bool,
) error {
	var (
		ctx        = convert(breaker)
		err  error = internal
		core error
	)

	for attempt, should := uint(0), true; should; attempt++ {
		core = unwrap(err)
		for i, repeat := 0, len(strategies); should && i < repeat; i++ {
			should = should && strategies[i](breaker, attempt, core)
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
	breaker Breaker,
	action func(context.Context) error,
	strategies ...func(Breaker, uint, error) bool,
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
