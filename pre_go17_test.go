// +build !go1.7

package retrier

import (
	"errors"
	"testing"
	"time"

	"github.com/kamilsk/retrier/strategy"
	"golang.org/x/net/context"
)

func TestRetry(t *testing.T) {
	action := func(attempt uint) error {
		return nil
	}

	err := Retry(context.Background(), action)

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

	err := Retry(context.Background(), action, strategy.Infinite())

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
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := Retry(ctx, func(uint) error { return nil }, strategy.Infinite()); err != ctx.Err() {
		t.Errorf("expected context done error, obtained %+v", err)
	}
}

func TestRetry_RetriesWithDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	action := func(uint) error {
		time.Sleep(110 * time.Millisecond)
		return nil
	}

	if err := Retry(ctx, action, strategy.Infinite()); err != ctx.Err() {
		t.Errorf("expected context done error, obtained %+v", err)
	}
}
