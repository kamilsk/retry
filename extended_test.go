package retry

import (
	"errors"
	"testing"

	"github.com/kamilsk/retry/strategy"
)

func TestRetryWithError(t *testing.T) {
	action := func(attempt uint) error {
		return nil
	}

	err := RetryWithError(action)

	if nil != err {
		t.Error("expected a nil error")
	}
}

func TestRetryWithErrorRetriesUntilNoErrorReturned(t *testing.T) {
	const errorUntilAttemptNumber = 5

	var attemptsMade uint

	action := func(attempt uint) error {
		attemptsMade = attempt

		if errorUntilAttemptNumber == attempt {
			return nil
		}

		return errors.New("erroring")
	}

	err := RetryWithError(action, strategy.ExtendStrategies(strategy.Infinite())...)

	if nil != err {
		t.Error("expected a nil error")
	}

	if errorUntilAttemptNumber != attemptsMade {
		t.Errorf(
			"expected %d attempts to be made, but %d were made instead",
			errorUntilAttemptNumber,
			attemptsMade,
		)
	}
}

func TestShouldAttemptWithError(t *testing.T) {
	var err error
	shouldAttempt := shouldAttemptWithError(1, err)

	if !shouldAttempt {
		t.Error("expected to return true")
	}
}

func TestShouldAttemptWithErrorAndStrategy(t *testing.T) {
	const attemptNumberShouldReturnFalse = 7
	var err error

	strategy := func(attempt uint, err error) bool {
		return (attemptNumberShouldReturnFalse != attempt)
	}

	should := shouldAttemptWithError(1, err, strategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttemptWithError(1+attemptNumberShouldReturnFalse, err, strategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttemptWithError(attemptNumberShouldReturnFalse, err, strategy)

	if should {
		t.Error("expected to return false")
	}
}

func TestShouldAttemptWithErrorAndMultipleStrategies(t *testing.T) {
	var err error
	trueStrategy := func(attempt uint, err error) bool {
		return true
	}

	falseStrategy := func(attempt uint, err error) bool {
		return false
	}

	should := shouldAttemptWithError(1, err, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttemptWithError(1, err, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttemptWithError(1, err, trueStrategy, trueStrategy, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttemptWithError(1, err, falseStrategy, falseStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttemptWithError(1, err, trueStrategy, trueStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}
}
