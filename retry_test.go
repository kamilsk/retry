package retry

import "testing"

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
