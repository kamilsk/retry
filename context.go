package retry

import "context"

// TryContext takes an interruptable action and performs it, repetitively, until successful.
// It uses the Context as a Breaker to prevent unnecessary action execution.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
//
// TODO:v5 not quite honest implementation
func TryContext(
	ctx context.Context,
	action func(ctx context.Context, attempt uint) error,
	strategies ...func(attempt uint, err error) bool,
) (err error) {

	// TODO:v5 will be removed
	defer func() {
		if r := recover(); r != nil {
			err = result{err, r}
		}
	}()

	for attempt, repeat, should := uint(0), len(strategies), true; should; attempt++ {
		for i := 0; should && i < repeat; i++ {
			should = should && ctx.Err() == nil && strategies[i](attempt, err)
		}

		if !should && ctx.Err() != nil {
			return Interrupted
		}

		if should {
			err = action(ctx, attempt)
			should = err != nil
		}
	}

	return err
}
