package retry

import "github.com/kamilsk/retry/strategy"

// Action defines a callable function that package retry can handle.
type Action func(attempt uint) error

// Retry takes an action and performs it, repetitively, until successful.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Retry(action Action, strategies ...strategy.Strategy) error {
	var err error
	attempt := uint(0)

	if len(strategies) == 0 {
		return action(attempt)
	}

	for ; (0 == attempt || nil != err) && shouldAttempt(attempt, strategies...); attempt++ {
		err = action(attempt)
	}

	return err
}

// shouldAttempt evaluates the provided strategies with the given attempt to
// determine if the Retry loop should make another attempt.
func shouldAttempt(attempt uint, strategies ...strategy.Strategy) bool {
	shouldAttempt := true

	for i := 0; shouldAttempt && i < len(strategies); i++ {
		shouldAttempt = shouldAttempt && strategies[i](attempt)
	}

	return shouldAttempt
}
