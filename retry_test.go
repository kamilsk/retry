package retry_test

import (
	"errors"
	"testing"
	"time"

	"github.com/kamilsk/retry"
	"github.com/kamilsk/retry/strategy"
)

func TestRetry(t *testing.T) {
	action := func(attempt uint) error {
		return nil
	}

	err := retry.Retry(nil, action)

	if nil != err {
		t.Error("expected a nil error")
	}
}

func TestRetry_RetriesUntilNoErrorReturned(t *testing.T) {
	const errorUntilAttemptNumber = 5

	var attemptsMade uint

	action := func(attempt uint) error {
		attemptsMade = attempt

		if errorUntilAttemptNumber == attempt {
			return nil
		}

		return errors.New("erroring")
	}

	err := retry.Retry(nil, action, strategy.Infinite())

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

func TestRetry_RetriesWithAlreadyDoneContext(t *testing.T) {
	deadline, expected := retry.WithTimeout(0), "operation timeout"

	if err := retry.Retry(deadline, func(uint) error { return nil }, strategy.Infinite()); !retry.IsTimeout(err) {
		t.Errorf("an unexpected error. expected: %s; obtained: %v", expected, err)
	}
}

func TestRetry_RetriesWithDeadline(t *testing.T) {
	deadline, expected := retry.WithTimeout(100*time.Millisecond), "operation timeout"

	action := func(uint) error {
		time.Sleep(110 * time.Millisecond)
		return nil
	}

	if err := retry.Retry(deadline, action, strategy.Infinite()); !retry.IsTimeout(err) {
		t.Errorf("an unexpected error. expected: %s; obtained: %v", expected, err)
	}
}
