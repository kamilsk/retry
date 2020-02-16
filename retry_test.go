package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/strategy"
)

func TestDo(t *testing.T) {
	type Assert struct {
		Attempts uint
		Error    func(error) bool
	}

	tests := map[string]struct {
		breaker    strategy.Breaker
		strategies How
		error      error
		assert     Assert
	}{
		"zero iterations": {
			newClosedBreaker(),
			How{
				strategy.Delay(10 * time.Millisecond),
				strategy.Limit(10000),
			},
			errors.New("zero iterations"),
			Assert{0, func(err error) bool { return err == context.Canceled }},
		},
		"one iteration": {
			newBreaker(),
			nil,
			nil,
			Assert{1, func(err error) bool { return err == nil }},
		},
		"two iterations": {
			newBreaker(),
			How{strategy.Limit(2)},
			errors.New("two iterations"),
			Assert{2, func(err error) bool { return err != nil && err.Error() == "two iterations" }},
		},
		"three iterations": {
			newBreaker(),
			How{strategy.Limit(3)},
			errors.New("three iterations"),
			Assert{3, func(err error) bool { return err != nil && err.Error() == "three iterations" }},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var total uint
			action := func() error {
				total += 1
				return test.error
			}
			err := Do(test.breaker, action, test.strategies...)
			if !test.assert.Error(err) {
				t.Error("fail error assertion")
			}
			if test.assert.Attempts != total {
				t.Errorf("expected %d attempts, obtained %d", test.assert.Attempts, total)
			}
		})
	}
}

func newBreaker() *contextBreaker {
	ctx, cancel := context.WithCancel(context.Background())
	return &contextBreaker{ctx, cancel}
}

func newClosedBreaker() *contextBreaker {
	breaker := newBreaker()
	breaker.Close()
	return breaker
}

type contextBreaker struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (breaker *contextBreaker) Done() <-chan struct{} {
	return breaker.ctx.Done()
}

func (breaker *contextBreaker) Close() {
	breaker.cancel()
}

func (breaker *contextBreaker) Err() error {
	return breaker.ctx.Err()
}
