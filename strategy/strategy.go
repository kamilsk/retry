// Package strategy provides a way to define how retry is performed.
package strategy

import (
	"net"
	"time"
)

// Strategy defines a function that Retry calls before every successive attempt
// to determine whether it should make the next attempt or not. Returning true
// allows for the next attempt to be made. Returning false halts the retrying
// process and returns the last error returned by the called Action.
//
// The strategy will be passed an "attempt" number on each successive retry
// iteration, starting with a 0 value before the first attempt is actually
// made. This allows for a pre-action delay, etc.
type Strategy func(breaker Breaker, attempt uint, err error) bool

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
		if attempt == 0 {
			timer := time.NewTimer(duration)
			select {
			case <-breaker.Done():
				_ = timer.Stop()
				return false
			case <-timer.C:
				_ = timer.Stop()
			}
		}
		return true
	}
}

// Wait creates a Strategy that waits the given durations for each attempt after
// the first. If the number of attempts is greater than the number of durations
// provided, then the strategy uses the last duration provided.
func Wait(durations ...time.Duration) Strategy {
	return func(breaker Breaker, attempt uint, _ error) bool {
		if attempt > 0 && len(durations) > 0 {
			durationIndex := int(attempt - 1)
			if len(durations) <= durationIndex {
				durationIndex = len(durations) - 1
			}
			timer := time.NewTimer(durations[durationIndex])
			select {
			case <-breaker.Done():
				_ = timer.Stop()
				return false
			case <-timer.C:
				_ = timer.Stop()
			}
		}
		return true
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
		if attempt > 0 {
			timer := time.NewTimer(transformation(algorithm(attempt)))
			select {
			case <-breaker.Done():
				_ = timer.Stop()
				return false
			case <-timer.C:
				_ = timer.Stop()
			}
		}
		return true
	}
}

const (
	Skip   = true
	Strict = false
)

// ErrorHandler defines a function that CheckError can use
// to determine whether it should make the next attempt or not.
// Returning true allows for the next attempt to be made.
// Returning false halts the retrying process and returns the last error
// returned by the called Action.
type ErrorHandler func(error) bool

// CheckError creates a Strategy that checks an error and returns
// if an error is retriable or not. Otherwise, it returns the defaults.
func CheckError(handlers ...func(error) bool) Strategy {
	// equal to go.octolab.org/errors.Retriable
	type retriable interface {
		error
		Retriable() bool // Is the error retriable?
	}

	return func(_ Breaker, _ uint, err error) bool {
		if err == nil {
			return true
		}
		if err, is := err.(retriable); is {
			return err.Retriable()
		}
		for _, handle := range handlers {
			if !handle(err) {
				return false
			}
		}
		return true
	}
}

// NetworkError creates an error Handler that checks an error and returns true
// if an error is the temporary network error.
// The Handler returns the defaults if an error is not a network error.
func NetworkError(defaults bool) func(error) bool {
	return func(err error) bool {
		if err, is := err.(net.Error); is {
			return err.Temporary() || err.Timeout()
		}
		return defaults
	}
}
