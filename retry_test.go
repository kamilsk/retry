package retry

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

var delta = 10 * time.Millisecond

func TestRetry(t *testing.T) {
	type Assert struct {
		Attempts uint
		Error    func(error) bool
	}

	tests := []struct {
		name       string
		breaker    BreakCloser
		strategies []func(attempt uint, err error) bool
		error      error
		assert     Assert
	}{
		{
			"zero iterations",
			newClosedBreaker(),
			[]func(attempt uint, err error) bool{delay(delta), limit(10000)},
			errors.New("zero iterations"),
			Assert{0, IsInterrupted},
		},
		{
			"one iteration",
			nil,
			nil,
			nil,
			Assert{1, func(err error) bool { return err == nil }},
		},
		{
			"two iterations",
			nil,
			[]func(attempt uint, err error) bool{limit(2)},
			errors.New("two iterations"),
			Assert{2, func(err error) bool { return err != nil }},
		},
		{
			"three iterations",
			newBreaker(),
			[]func(attempt uint, err error) bool{limit(3)},
			errors.New("three iterations"),
			Assert{3, func(err error) bool { return err != nil && err.Error() == "three iterations" }},
		},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			var total uint
			action := func(attempt uint) error {
				total = attempt + 1
				return tc.error
			}
			err := Retry(tc.breaker, action, tc.strategies...)
			if !tc.assert.Error(err) {
				t.Error("fail error assertion")
			}
			if _, is := IsRecovered(err); is {
				t.Error("recovered panic is not expected")
			}
			if tc.assert.Attempts != total {
				t.Errorf("expected %d attempts, obtained %d", tc.assert.Attempts, total)
			}
		})
	}
	t.Run("unexpected panic", func(t *testing.T) {
		err := Retry(newBreaker(), func(uint) error { panic("Catch Me If You Can") })
		cause, is := IsRecovered(err)
		if !is {
			t.Fatal("recovered panic is expected")
		}
		if !reflect.DeepEqual(cause, "Catch Me If You Can") {
			t.Fatal("Catch Me If You Can is expected")
		}
	})
}

func TestTry(t *testing.T) {
	type Assert struct {
		Attempts uint
		Error    func(error) bool
	}

	tests := []struct {
		name       string
		breaker    Breaker
		strategies []func(attempt uint, err error) bool
		error      error
		assert     Assert
	}{
		{
			"zero iterations",
			newClosedBreaker(),
			[]func(attempt uint, err error) bool{delay(delta), limit(10000)},
			errors.New("zero iterations"),
			Assert{0, IsInterrupted},
		},
		{
			"one iteration",
			nil,
			nil,
			nil,
			Assert{1, func(err error) bool { return err == nil }},
		},
		{
			"two iterations",
			nil,
			[]func(attempt uint, err error) bool{limit(2)},
			errors.New("two iterations"),
			Assert{2, func(err error) bool { return err != nil }},
		},
		{
			"three iterations",
			newPanicBreaker(),
			[]func(attempt uint, err error) bool{limit(3)},
			errors.New("three iterations"),
			Assert{3, func(err error) bool { return err != nil && err.Error() == "three iterations" }},
		},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			var total uint
			action := func(attempt uint) error {
				total = attempt + 1
				return tc.error
			}
			err := Try(tc.breaker, action, tc.strategies...)
			if !tc.assert.Error(err) {
				t.Error("fail error assertion")
			}
			if _, is := IsRecovered(err); is {
				t.Error("recovered panic is not expected")
			}
			if tc.assert.Attempts != total {
				t.Errorf("expected %d attempts, obtained %d", tc.assert.Attempts, total)
			}
		})
	}
	t.Run("unexpected panic", func(t *testing.T) {
		err := Try(newBreaker(), func(uint) error { panic("Catch Me If You Can") })
		cause, is := IsRecovered(err)
		if !is {
			t.Fatal("recovered panic is expected")
		}
		if !reflect.DeepEqual(cause, "Catch Me If You Can") {
			t.Fatal("Catch Me If You Can is expected")
		}
	})
}

func TestTryContext(t *testing.T) {
	type Assert struct {
		Attempts uint
		Error    func(error) bool
	}

	tests := []struct {
		name       string
		ctx        context.Context
		strategies []func(attempt uint, err error) bool
		error      error
		assert     Assert
	}{
		{
			"zero iterations",
			func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			[]func(attempt uint, err error) bool{delay(delta), limit(10000)},
			errors.New("zero iterations"),
			Assert{0, IsInterrupted},
		},
		{
			"one iteration",
			nil,
			nil,
			nil,
			Assert{1, func(err error) bool { return err == nil }},
		},
		{
			"two iterations",
			nil,
			[]func(attempt uint, err error) bool{limit(2)},
			errors.New("two iterations"),
			Assert{2, func(err error) bool { return err != nil }},
		},
		{
			"three iterations",
			context.Background(),
			[]func(attempt uint, err error) bool{limit(3)},
			errors.New("three iterations"),
			Assert{3, func(err error) bool { return err != nil && err.Error() == "three iterations" }},
		},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			var total uint
			action := func(ctx context.Context, attempt uint) error {
				if !reflect.DeepEqual(tc.ctx, ctx) {
					t.Fatal("an unexpected context obtained")
				}
				total = attempt + 1
				return tc.error
			}
			err := TryContext(tc.ctx, action, tc.strategies...)
			if !tc.assert.Error(err) {
				t.Error("fail error assertion")
			}
			if _, is := IsRecovered(err); is {
				t.Error("recovered panic is not expected")
			}
			if tc.assert.Attempts != total {
				t.Errorf("expected %d attempts, obtained %d", tc.assert.Attempts, total)
			}
		})
	}
	t.Run("unexpected panic", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		err := TryContext(ctx, func(context.Context, uint) error { panic("Catch Me If You Can") })
		cause, is := IsRecovered(err)
		if !is {
			t.Fatal("recovered panic is expected")
		}
		if !reflect.DeepEqual(cause, "Catch Me If You Can") {
			t.Fatal("Catch Me If You Can is expected")
		}
		cancel()
	})
}
