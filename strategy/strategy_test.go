package strategy_test

import (
	"context"
	"testing"
	"time"

	. "github.com/kamilsk/retry/v5/strategy"
)

// timeMarginOfError represents the acceptable amount of time that may pass for
// a time-based (sleep) unit before considering invalid.
const timeMarginOfError = time.Millisecond

func TestLimit(t *testing.T) {
	breaker, policy := context.Background(), Limit(3)

	if !policy(breaker, 0) {
		t.Error("strategy expected to return true")
	}

	if !policy(breaker, 1) {
		t.Error("strategy expected to return true")
	}

	if !policy(breaker, 2) {
		t.Error("strategy expected to return true")
	}

	if policy(breaker, 3) {
		t.Error("strategy expected to return false")
	}
}

func TestDelay(t *testing.T) {
	const delayDuration = 10 * timeMarginOfError

	breaker, policy := context.Background(), Delay(delayDuration)

	if now := time.Now(); !policy(breaker, 0) || delayDuration > time.Since(now) {
		t.Errorf("strategy expected to return true in %s", delayDuration)
	}

	if now := time.Now(); !policy(breaker, 5) || (delayDuration/10) < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}
}

func TestWait(t *testing.T) {
	breaker, policy := context.Background(), Wait()

	if now := time.Now(); !policy(breaker, 0) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	if now := time.Now(); !policy(breaker, 999) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}
}

func TestWaitWithDuration(t *testing.T) {
	const waitDuration = 10 * timeMarginOfError

	breaker, policy := context.Background(), Wait(waitDuration)

	if now := time.Now(); !policy(breaker, 0) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	if now := time.Now(); !policy(breaker, 1) || waitDuration > time.Since(now) {
		t.Errorf("strategy expected to return true in %s", waitDuration)
	}
}

func TestWaitWithMultipleDurations(t *testing.T) {
	waitDurations := []time.Duration{
		10 * timeMarginOfError,
		20 * timeMarginOfError,
		30 * timeMarginOfError,
		40 * timeMarginOfError,
	}

	breaker, policy := context.Background(), Wait(waitDurations...)

	if now := time.Now(); !policy(breaker, 0) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	if now := time.Now(); !policy(breaker, 1) || waitDurations[0] > time.Since(now) {
		t.Errorf("strategy expected to return true in %s", waitDurations[0])
	}

	if now := time.Now(); !policy(breaker, 3) || waitDurations[2] > time.Since(now) {
		t.Errorf("strategy expected to return true in %s", waitDurations[2])
	}

	if now := time.Now(); !policy(breaker, 999) || waitDurations[len(waitDurations)-1] > time.Since(now) {
		t.Errorf("strategy expected to return true in %s", waitDurations[len(waitDurations)-1])
	}
}

func TestBackoff(t *testing.T) {
	const backoffDuration = 10 * timeMarginOfError
	const algorithmDurationBase = timeMarginOfError

	algorithm := func(attempt uint) time.Duration {
		return backoffDuration - (algorithmDurationBase * time.Duration(attempt))
	}

	breaker, policy := context.Background(), Backoff(algorithm)

	if now := time.Now(); !policy(breaker, 0) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	for i := uint(1); i < 10; i++ {
		expectedResult := algorithm(i)

		if now := time.Now(); !policy(breaker, i) || expectedResult > time.Since(now) {
			t.Errorf("strategy expected to return true in %s", expectedResult)
		}
	}
}

func TestBackoffWithJitter(t *testing.T) {
	const backoffDuration = 10 * timeMarginOfError
	const algorithmDurationBase = timeMarginOfError

	algorithm := func(attempt uint) time.Duration {
		return backoffDuration - (algorithmDurationBase * time.Duration(attempt))
	}

	transformation := func(duration time.Duration) time.Duration {
		return duration - 10*timeMarginOfError
	}

	breaker, policy := context.Background(), BackoffWithJitter(algorithm, transformation)

	if now := time.Now(); !policy(breaker, 0) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	for i := uint(1); i < 10; i++ {
		expectedResult := transformation(algorithm(i))

		if now := time.Now(); !policy(breaker, i) || expectedResult > time.Since(now) {
			t.Errorf("strategy expected to return true in %s", expectedResult)
		}
	}
}
