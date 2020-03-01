package strategy_test

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"testing"
	"time"

	. "github.com/kamilsk/retry/v5/strategy"
)

func TestLimit(t *testing.T) {
	tests := map[string]struct {
		value    uint
		args     tuple
		expected bool
	}{
		"first call": {
			2,
			tuple{breaker(), 0, nil},
			true,
		},
		"first call with error": {
			2,
			tuple{breaker(), 0, errors.New("ignored")},
			true,
		},
		"first call with interrupted breaker": {
			2,
			tuple{interrupted(), 0, nil},
			true,
		},
		"next call": {
			2,
			tuple{breaker(), 1, nil},
			true,
		},
		"next call with error": {
			2,
			tuple{breaker(), 1, errors.New("ignored")},
			true,
		},
		"next call with interrupted breaker": {
			2,
			tuple{interrupted(), 1, nil},
			true,
		},
		"last call": {
			2,
			tuple{breaker(), 999, nil},
			false,
		},
		"last call with error": {
			2,
			tuple{breaker(), 999, errors.New("ignored")},
			false,
		},
		"last call with interrupted breaker": {
			2,
			tuple{interrupted(), 999, nil},
			false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			policy := Limit(test.value)
			if obtained := policy(test.args.unpack()); test.expected != obtained {
				t.Errorf("expected: %v, obtained: %v", test.expected, obtained)
			}
		})
	}
}

func TestDelay(t *testing.T) {
	tests := map[string]struct {
		duration time.Duration
		args     tuple
		expected bool
		assert   func(time.Time, time.Duration) bool
	}{
		"first call": {
			time.Millisecond,
			tuple{breaker(), 0, nil},
			true,
			func(past time.Time, expected time.Duration) bool {
				return time.Since(past) > expected
			},
		},
		"first call with error": {
			time.Millisecond,
			tuple{breaker(), 0, errors.New("ignored")},
			true,
			func(past time.Time, expected time.Duration) bool {
				return time.Since(past) > expected
			},
		},
		"first call with interrupted breaker": {
			time.Millisecond,
			tuple{interrupted(), 0, nil},
			false,
			func(past time.Time, expected time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"next call": {
			time.Millisecond,
			tuple{breaker(), 999, nil},
			true,
			func(past time.Time, expected time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"next call with error": {
			time.Millisecond,
			tuple{breaker(), 999, errors.New("ignored")},
			true,
			func(past time.Time, expected time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"next call with interrupted breaker": {
			time.Millisecond,
			tuple{interrupted(), 999, nil},
			true,
			func(past time.Time, expected time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			policy, now := Delay(test.duration), time.Now()
			if obtained := policy(test.args.unpack()); test.expected != obtained {
				t.Errorf("expected: %v, obtained: %v", test.expected, obtained)
			}
			if !test.assert(now, test.duration) {
				t.Error("delay is not asserted")
			}
		})
	}
}

func TestWait(t *testing.T) {
	tests := map[string]struct {
		durations []time.Duration
		args      tuple
		expected  bool
		assert    func(time.Time, []time.Duration) bool
	}{
		"first call with empty durations": {
			nil,
			tuple{breaker(), 0, nil},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"first call with empty durations and error": {
			nil,
			tuple{breaker(), 0, errors.New("ignored")},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"first call with empty durations and interrupted breaker": {
			nil,
			tuple{interrupted(), 0, nil},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"next call with empty durations": {
			nil,
			tuple{breaker(), 999, nil},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"next call with empty durations and error": {
			nil,
			tuple{breaker(), 999, errors.New("ignored")},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"next call with empty durations and interrupted breaker": {
			nil,
			tuple{interrupted(), 999, nil},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"first call with multiple durations": {
			[]time.Duration{time.Minute, time.Hour},
			tuple{breaker(), 0, nil},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"first call with multiple durations and error": {
			[]time.Duration{time.Minute, time.Hour},
			tuple{breaker(), 0, errors.New("ignored")},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"first call with multiple durations and interrupted breaker": {
			[]time.Duration{time.Minute, time.Hour},
			tuple{interrupted(), 0, nil},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"next call with multiple durations": {
			[]time.Duration{time.Millisecond, time.Hour},
			tuple{breaker(), 1, nil},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) > durations[0]
			},
		},
		"next call with multiple durations and error": {
			[]time.Duration{time.Millisecond, time.Hour},
			tuple{breaker(), 1, errors.New("ignored")},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) > durations[0]
			},
		},
		"next call with multiple durations and interrupted breaker": {
			[]time.Duration{time.Minute, time.Hour},
			tuple{interrupted(), 1, nil},
			false,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"last call with multiple durations": {
			[]time.Duration{time.Hour, time.Millisecond},
			tuple{breaker(), 999, nil},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) > durations[len(durations)-1]
			},
		},
		"last call with multiple durations and error": {
			[]time.Duration{time.Hour, time.Millisecond},
			tuple{breaker(), 999, errors.New("ignored")},
			true,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) > durations[len(durations)-1]
			},
		},
		"last call with multiple durations and interrupted breaker": {
			[]time.Duration{time.Minute, time.Hour},
			tuple{interrupted(), 999, nil},
			false,
			func(past time.Time, durations []time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			policy, now := Wait(test.durations...), time.Now()
			if obtained := policy(test.args.unpack()); test.expected != obtained {
				t.Errorf("expected: %v, obtained: %v", test.expected, obtained)
			}
			if !test.assert(now, test.durations) {
				t.Error("wait is not asserted")
			}
		})
	}
}

func TestBackoff(t *testing.T) {
	tests := map[string]struct {
		algorithm func(uint) time.Duration
		args      tuple
		expected  bool
		assert    func(time.Time, time.Duration) bool
	}{
		"first call": {
			func(uint) time.Duration { return time.Hour },
			tuple{breaker(), 0, nil},
			true,
			func(past time.Time, duration time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"first call with error": {
			func(uint) time.Duration { return time.Hour },
			tuple{breaker(), 0, errors.New("ignored")},
			true,
			func(past time.Time, duration time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"first call with interrupted breaker": {
			func(uint) time.Duration { return time.Hour },
			tuple{interrupted(), 0, nil},
			true,
			func(past time.Time, duration time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"next call": {
			func(attempt uint) time.Duration {
				return time.Duration(attempt) * time.Millisecond
			},
			tuple{breaker(), 2, nil},
			true,
			func(past time.Time, duration time.Duration) bool {
				return time.Since(past) > duration
			},
		},
		"next call with error": {
			func(attempt uint) time.Duration {
				return time.Duration(attempt) * time.Millisecond
			},
			tuple{breaker(), 2, errors.New("ignored")},
			true,
			func(past time.Time, duration time.Duration) bool {
				return time.Since(past) > duration
			},
		},
		"next call with interrupted breaker": {
			func(uint) time.Duration { return time.Hour },
			tuple{interrupted(), 999, nil},
			false,
			func(past time.Time, duration time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			policy, now := Backoff(test.algorithm), time.Now()
			if obtained := policy(test.args.unpack()); test.expected != obtained {
				t.Errorf("expected: %v, obtained: %v", test.expected, obtained)
			}
			if !test.assert(now, test.algorithm(test.args.attempt)) {
				t.Error("backoff is not asserted")
			}
		})
	}
}

func TestBackoffWithJitter(t *testing.T) {
	tests := map[string]struct {
		algorithm      func(uint) time.Duration
		transformation func(time.Duration) time.Duration
		args           tuple
		expected       bool
		assert         func(time.Time, time.Duration) bool
	}{
		"first call": {
			func(uint) time.Duration { return time.Hour },
			func(duration time.Duration) time.Duration { return duration },
			tuple{breaker(), 0, nil},
			true,
			func(past time.Time, duration time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"first call with error": {
			func(uint) time.Duration { return time.Hour },
			func(duration time.Duration) time.Duration { return duration },
			tuple{breaker(), 0, errors.New("ignored")},
			true,
			func(past time.Time, duration time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"first call with interrupted breaker": {
			func(uint) time.Duration { return time.Hour },
			func(duration time.Duration) time.Duration { return duration },
			tuple{interrupted(), 0, nil},
			true,
			func(past time.Time, duration time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
		"next call": {
			func(attempt uint) time.Duration {
				return time.Hour + time.Duration(attempt)*time.Millisecond
			},
			func(duration time.Duration) time.Duration {
				return duration - time.Hour
			},
			tuple{breaker(), 2, nil},
			true,
			func(past time.Time, duration time.Duration) bool {
				return time.Since(past) > duration
			},
		},
		"next call with error": {
			func(attempt uint) time.Duration {
				return time.Hour + time.Duration(attempt)*time.Millisecond
			},
			func(duration time.Duration) time.Duration {
				return duration - time.Hour
			},
			tuple{breaker(), 2, errors.New("ignored")},
			true,
			func(past time.Time, duration time.Duration) bool {
				return time.Since(past) > duration
			},
		},
		"next call with interrupted breaker": {
			func(uint) time.Duration { return time.Hour },
			func(duration time.Duration) time.Duration { return duration },
			tuple{interrupted(), 999, nil},
			false,
			func(past time.Time, duration time.Duration) bool {
				return time.Since(past) < time.Millisecond
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			policy, now := BackoffWithJitter(test.algorithm, test.transformation), time.Now()
			if obtained := policy(test.args.unpack()); test.expected != obtained {
				t.Errorf("expected: %v, obtained: %v", test.expected, obtained)
			}
			if !test.assert(now, test.transformation(test.algorithm(test.args.attempt))) {
				t.Error("backoff with jitter is not asserted")
			}
		})
	}
}

func TestCheckError(t *testing.T) {
	generator := rand.New(rand.NewSource(0))

	tests := map[string]struct {
		error    error
		handlers []func(error) bool
		expected bool
	}{
		"nil error without handlers": {
			nil,
			nil,
			true,
		},
		"nil error with strict handlers": {
			nil,
			[]func(error) bool{NetworkError(Strict)},
			true,
		},
		"retriable error without handlers": {
			retriable("yes"),
			nil,
			true,
		},
		"retriable error with strict handlers": {
			retriable("yes"),
			[]func(error) bool{NetworkError(Strict)},
			true,
		},
		"non-retriable error without handlers": {
			retriable("no"),
			nil,
			false,
		},
		"non-retriable error with strict handlers": {
			retriable("no"),
			nil,
			false,
		},
		"network address error with strict check": {
			&net.AddrError{},
			[]func(error) bool{NetworkError(Strict)},
			false,
		},
		"network address error without strict check": {
			&net.AddrError{},
			[]func(error) bool{NetworkError(Skip)},
			false,
		},
		"temporary dns error": {
			&net.DNSError{IsTemporary: true},
			[]func(error) bool{NetworkError(Strict)},
			true,
		},
		"an error with strict check": {
			errors.New("test"),
			[]func(error) bool{NetworkError(Strict)},
			false,
		},
		"an error without strict check": {
			errors.New("test"),
			[]func(error) bool{NetworkError(Skip)},
			true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			policy, attempt := CheckError(test.handlers...), uint(generator.Uint32())
			if obtained := policy(breaker(), attempt, test.error); test.expected != obtained {
				t.Errorf("expected: %v, obtained: %v", test.expected, obtained)
			}
		})
	}
}

// helpers

func breaker() Breaker {
	return context.Background()
}

func interrupted() Breaker {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

type tuple struct {
	breaker Breaker
	attempt uint
	error   error
}

func (tuple *tuple) unpack() (Breaker, uint, error) {
	return tuple.breaker, tuple.attempt, tuple.error
}

type retriable string

func (err retriable) Error() string   { return string(err) }
func (err retriable) Retriable() bool { return err == "yes" }
