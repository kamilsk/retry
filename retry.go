// Package retrier provides functional mechanism based on context
// to perform actions repetitively until successful.
//
// Copyright © 2016 Trevor N. Suarez (Rican7)
package retrier

import "github.com/kamilsk/retrier/strategy"

// Action defines a callable function that package retrier can handle.
type Action func(attempt uint) error

// shouldAttempt evaluates the provided strategies with the given attempt to
// determine if the Retry loop should make another attempt.
func shouldAttempt(attempt uint, err error, strategies ...strategy.Strategy) bool {
	shouldAttempt := true

	for i, repeat := 0, len(strategies); shouldAttempt && i < repeat; i++ {
		shouldAttempt = shouldAttempt && strategies[i](attempt, err)
	}

	return shouldAttempt
}
