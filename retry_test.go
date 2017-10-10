package retry

import (
	"errors"
	"testing"

	"github.com/kamilsk/retry/strategy"
)

func TestRetry(t *testing.T) {
	action := func(attempt uint) error {
		return nil
	}

	err := Retry(nil, action)

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

	err := Retry(nil, action, strategy.Infinite())

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

// TODO ctx.Err() should be replaced correctly, atomic is good candidate
//func TestRetry_RetriesWithAlreadyDoneContext(t *testing.T) {
//	deadline := WithTimeout(0)
//
//	if err := Retry(deadline, func(uint) error { return nil }, strategy.Infinite()); err != ctx.Err() {
//		t.Errorf("expected context done error, obtained %+v", err)
//	}
//}

// TODO ctx.Err() should be replaced correctly, atomic is good candidate
//func TestRetry_RetriesWithDeadline(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
//	defer cancel()
//
//	action := func(uint) error {
//		time.Sleep(110 * time.Millisecond)
//		return nil
//	}
//
//	if err := Retry(ctx, action, strategy.Infinite()); err != ctx.Err() {
//		t.Errorf("expected context done error, obtained %+v", err)
//	}
//}

func TestShouldAttempt(t *testing.T) {
	shouldAttempt := shouldAttempt(1, nil)

	if !shouldAttempt {
		t.Error("expected to return true")
	}
}

func TestShouldAttemptWithStrategy(t *testing.T) {
	const attemptNumberShouldReturnFalse = 7

	s := func(attempt uint, _ error) bool {
		return attemptNumberShouldReturnFalse != attempt
	}

	should := shouldAttempt(1, nil, s)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1+attemptNumberShouldReturnFalse, nil, s)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(attemptNumberShouldReturnFalse, nil, s)

	if should {
		t.Error("expected to return false")
	}
}

func TestShouldAttemptWithMultipleStrategies(t *testing.T) {
	trueStrategy := func(attempt uint, _ error) bool {
		return true
	}

	falseStrategy := func(attempt uint, _ error) bool {
		return false
	}

	should := shouldAttempt(1, nil, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1, nil, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttempt(1, nil, trueStrategy, trueStrategy, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1, nil, falseStrategy, falseStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttempt(1, nil, trueStrategy, trueStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}
}
