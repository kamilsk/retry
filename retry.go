// Package retry provides the most advanced interruptible mechanism
// to perform actions repetitively until successful.
package retry

import "github.com/kamilsk/retry/v5/strategy"

// Action defines a callable function that package retry can handle.
type Action func() error

// How is an alias for batch of Strategies.
//
//  how := retry.How{
//  	strategy.Limit(3),
//  }
//
type How []func(strategy.Breaker, uint) bool

// Do takes an action and performs it, repetitively, until successful.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Do(
	breaker strategy.Breaker,
	action func() error,
	strategies ...func(strategy.Breaker, uint) bool,
) error {
	var err error
	done := make(chan struct{})

	go func() {
		for attempt := uint(0); shouldAttempt(breaker, attempt, err, strategies...); attempt++ {
			err = action()
		}
		close(done)
	}()

	select {
	case <-breaker.Done():
		return breaker.Err()
	case <-done:
		return err
	}
}

// shouldAttempt evaluates the provided strategies with the given attempt to
// determine if the Retry loop should make another attempt.
func shouldAttempt(breaker strategy.Breaker, attempt uint, err error, strategies ...func(strategy.Breaker, uint) bool) bool {
	should := attempt == 0 || err != nil

	for i, repeat := 0, len(strategies); should && i < repeat; i++ {
		should = should && strategies[i](breaker, attempt)
	}

	return should
}
