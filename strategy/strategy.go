// Package strategy provides a way to change the way that retry is performed.
package strategy

import (
	"time"

	"github.com/kamilsk/retry/backoff"
	"github.com/kamilsk/retry/jitter"
)

// Strategy defines a function that Retry calls before every successive attempt
// to determine whether it should make the next attempt or not. Returning `true`
// allows for the next attempt to be made. Returning `false` halts the retrying
// process and returns the last error returned by the called Action.
//
// The strategy will be passed an "attempt" number on each successive retry
// iteration, starting with a `0` value before the first attempt is actually
// made. This allows for a pre-action delay, etc.
type Strategy func(attempt uint) bool

// Infinite creates a Strategy that will never stop repeating.
func Infinite() Strategy {
	return func(attempt uint) bool {
		return true
	}
}

// Timeout creates a Strategy that will return false if time is over.
// Not thread-safe. Do not use the same instance in multiple goroutines simultaneously.
//
//  // The example below shows how to get a race condition.
//  func TestTimeout_Concurrently(t *testing.T) {
//  	strategy := Timeout(100 * time.Millisecond)
//
//  	start := make(chan bool)
//  	wg := &sync.WaitGroup{}
//
//  	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
//  		wg.Add(1)
//  		go func() {
//  			defer wg.Done()
//  			<-start
//  			strategy(0)
//  		}()
//  	}
//
//  	close(start)
//  	wg.Wait()
//  }
func Timeout(timeout time.Duration) Strategy {
	var start time.Time
	return func(attempt uint) bool {
		if attempt == 0 {
			start = time.Now()
			return true
		}
		return start.Add(timeout).After(time.Now())
	}
}

// Limit creates a Strategy that limits the number of attempts that Retry will
// make.
func Limit(attemptLimit uint) Strategy {
	return func(attempt uint) bool {
		return (attempt <= attemptLimit)
	}
}

// Delay creates a Strategy that waits the given duration before the first
// attempt is made.
func Delay(duration time.Duration) Strategy {
	return func(attempt uint) bool {
		if 0 == attempt {
			time.Sleep(duration)
		}

		return true
	}
}

// Wait creates a Strategy that waits the given durations for each attempt after
// the first. If the number of attempts is greater than the number of durations
// provided, then the strategy uses the last duration provided.
func Wait(durations ...time.Duration) Strategy {
	return func(attempt uint) bool {
		if 0 < attempt && 0 < len(durations) {
			durationIndex := int(attempt - 1)

			if len(durations) <= durationIndex {
				durationIndex = len(durations) - 1
			}

			time.Sleep(durations[durationIndex])
		}

		return true
	}
}

// Backoff creates a Strategy that waits before each attempt, with a duration as
// defined by the given backoff.Algorithm.
func Backoff(algorithm backoff.Algorithm) Strategy {
	return BackoffWithJitter(algorithm, noJitter())
}

// BackoffWithJitter creates a Strategy that waits before each attempt, with a
// duration as defined by the given backoff.Algorithm and jitter.Transformation.
func BackoffWithJitter(algorithm backoff.Algorithm, transformation jitter.Transformation) Strategy {
	return func(attempt uint) bool {
		if 0 < attempt {
			time.Sleep(transformation(algorithm(attempt)))
		}

		return true
	}
}

// noJitter creates a jitter.Transformation that simply returns the input.
func noJitter() jitter.Transformation {
	return func(duration time.Duration) time.Duration {
		return duration
	}
}
