// Package strategy provides a way to define how retry is performed.
package strategy

import "time"

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

// Strategy defines a function that Retry calls before every successive attempt
// to determine whether it should make the next attempt or not. Returning true
// allows for the next attempt to be made. Returning false halts the retrying
// process and returns the last error returned by the called Action.
//
// The strategy will be passed an "attempt" number on each successive retry
// iteration, starting with a 0 value before the first attempt is actually
// made. This allows for a pre-action delay, etc.
type Strategy = func(breaker Breaker, attempt uint, err error) bool

// Limit creates a Strategy that limits the number of attempts
// that Retry will make.
func Limit(value uint) Strategy {
	return func(_ Breaker, attempt uint, _ error) bool {
		return attempt < value
	}
}

// Delay creates a Strategy that waits the given duration
// before the first attempt is made.
func Delay(duration time.Duration) Strategy {
	return func(breaker Breaker, attempt uint, _ error) bool {
		keep := true
		if attempt == 0 {
			timer := time.NewTimer(duration)
			select {
			case <-timer.C:
			case <-breaker.Done():
				keep = false
			}
			stop(timer)
		}
		return keep
	}
}

// Wait creates a Strategy that waits the given durations for each attempt after
// the first. If the number of attempts is greater than the number of durations
// provided, then the strategy uses the last duration provided.
func Wait(durations ...time.Duration) Strategy {
	return func(breaker Breaker, attempt uint, _ error) bool {
		keep := true
		if attempt > 0 && len(durations) > 0 {
			durationIndex := int(attempt - 1)
			if len(durations) <= durationIndex {
				durationIndex = len(durations) - 1
			}
			timer := time.NewTimer(durations[durationIndex])
			select {
			case <-timer.C:
			case <-breaker.Done():
				keep = false
			}
			stop(timer)
		}
		return keep
	}
}

// Backoff creates a Strategy that waits before each attempt, with a duration as
// defined by the given backoff.Algorithm.
func Backoff(algorithm func(attempt uint) time.Duration) Strategy {
	return BackoffWithJitter(algorithm, func(duration time.Duration) time.Duration {
		return duration
	})
}

// BackoffWithJitter creates a Strategy that waits before each attempt, with a
// duration as defined by the given backoff.Algorithm and jitter.Transformation.
func BackoffWithJitter(
	algorithm func(attempt uint) time.Duration,
	transformation func(duration time.Duration) time.Duration,
) Strategy {
	return func(breaker Breaker, attempt uint, _ error) bool {
		keep := true
		if attempt > 0 {
			timer := time.NewTimer(transformation(algorithm(attempt)))
			select {
			case <-timer.C:
			case <-breaker.Done():
				keep = false
			}
			stop(timer)
		}
		return keep
	}
}

func stop(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
}
