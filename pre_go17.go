// +build !go1.7

package retrier

import (
	"github.com/kamilsk/retrier/strategy"
	"golang.org/x/net/context"
)

// Retry takes an action and performs it, repetitively, until successful.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Retry(ctx context.Context, action Action, strategies ...strategy.Strategy) error {
	var attempt uint

	if ctx.Err() != nil {
		return ctx.Err()
	}

	if len(strategies) == 0 {
		return action(attempt)
	}

	var err error
	done := make(chan struct{})
	go func() {
		for ; (attempt == 0 || err != nil) && shouldAttempt(attempt, err, strategies...) && ctx.Err() == nil; attempt++ {
			err = action(attempt)
		}
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return err
	}
}
