package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTryContext(t *testing.T) {
	type Assert struct {
		Attempts uint
		Error    func(error) bool
	}

	tests := map[string]struct {
		ctx        func() context.Context
		strategies []func(attempt uint, err error) bool
		error      error
		assert     Assert
	}{
		"zero iterations": {
			func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			[]func(attempt uint, err error) bool{delay(delta), limit(10000)},
			errors.New("zero iterations"),
			Assert{0, func(err error) bool { return err == Interrupted }},
		},
		"one iteration": {
			context.Background,
			nil,
			nil,
			Assert{1, func(err error) bool { return err == nil }},
		},
		"two iterations": {
			context.Background,
			[]func(attempt uint, err error) bool{limit(2)},
			errors.New("two iterations"),
			Assert{2, func(err error) bool { return err != nil && err.Error() == "two iterations" }},
		},
		"three iterations": {
			context.Background,
			[]func(attempt uint, err error) bool{limit(3)},
			errors.New("three iterations"),
			Assert{3, func(err error) bool { return err != nil && err.Error() == "three iterations" }},
		},
		"long iteration": {
			func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), delta)
				defer func() {
					time.Sleep(2 * delta)
					cancel()
				}()
				return ctx
			},
			[]func(attempt uint, err error) bool{delay(time.Hour)},
			errors.New("long iteration"),
			Assert{0, func(err error) bool { return err == Interrupted }},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var total uint
			action := func(ctx context.Context, attempt uint) error {
				total = attempt + 1
				return test.error
			}
			err := TryContext(test.ctx(), action, test.strategies...)
			if !test.assert.Error(err) {
				t.Error("fail error assertion")
			}
			if test.assert.Attempts != total {
				t.Errorf("expected %d attempts, obtained %d", test.assert.Attempts, total)
			}
		})
	}
}
