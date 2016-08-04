package retry

import (
	"github.com/kamilsk/retry/strategy"
)

// RetryWithError forwarding errors between repetitions.
func RetryWithError(action Action, strategies ...strategy.ExtendedStrategy) error {
	var err error
	attempt := uint(0)

	if len(strategies) == 0 {
		return action(attempt)
	}

	for ; (0 == attempt || nil != err) && shouldAttemptWithError(attempt, err, strategies...); attempt++ {
		err = action(attempt)
	}

	return err
}

func shouldAttemptWithError(attempt uint, err error, strategies ...strategy.ExtendedStrategy) bool {
	shouldAttempt := true

	for i := 0; shouldAttempt && i < len(strategies); i++ {
		shouldAttempt = shouldAttempt && strategies[i](attempt, err)
	}

	return shouldAttempt
}
