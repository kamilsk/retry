package retry_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/strategy"
)

func TestDo(t *testing.T) {
	tests := testCases

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var attempts uint
			action := func(ctx context.Context) error {
				attempts++
				return test.action(ctx)
			}
			err := Do(test.breaker, action, test.strategies...)
			if test.expected.attempts != attempts {
				t.Errorf("expected: %d, obtained: %d", test.expected.attempts, attempts)
			}
			if !reflect.DeepEqual(test.expected.error, err) {
				t.Errorf("expected: %#v, obtained: %#v", test.expected.error, err)
			}
		})
	}

	t.Run("preserve context values", func(t *testing.T) {
		val := "value"
		ctx := context.WithValue(context.TODO(), key{}, val)
		action := func(ctx context.Context) error {
			if !reflect.DeepEqual(val, ctx.Value(key{})) {
				t.Error("value is not preserved")
			}
			return nil
		}
		if err := Do(ctx, action); err != nil {
			t.Error("unexpected error")
		}
	})
}

func TestGo(t *testing.T) {
	tests := append(
		testCases,
		testCase{
			"action call with error panic",
			breaker(),
			How{strategy.Wait(time.Hour)},
			func(context.Context) error { panic(Error("failure")) },
			expected{1, Error("failure")},
		},
		testCase{
			"action call with non-error panic",
			breaker(),
			How{strategy.Wait(time.Hour)},
			func(context.Context) error { panic("non-error") },
			expected{1, fmt.Errorf("retry: unexpected panic: %#v", "non-error")},
		},
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var attempts uint
			action := func(ctx context.Context) error {
				attempts++
				return test.action(ctx)
			}
			err := Go(test.breaker, action, test.strategies...)
			if test.expected.attempts != attempts {
				t.Errorf("expected: %d, obtained: %d", test.expected.attempts, attempts)
			}
			if !reflect.DeepEqual(test.expected.error, err) {
				t.Errorf("expected: %#v, obtained: %#v", test.expected.error, err)
			}
		})
	}

	t.Run("preserve context values", func(t *testing.T) {
		val := "value"
		ctx := context.WithValue(context.TODO(), key{}, val)
		action := func(ctx context.Context) error {
			if !reflect.DeepEqual(val, ctx.Value(key{})) {
				t.Error("value is not preserved")
			}
			return nil
		}
		if err := Go(ctx, action); err != nil {
			t.Error("unexpected error")
		}
	})
}

// helpers

func breaker() Breaker {
	return context.TODO()
}

func interrupted() Breaker {
	ctx, cancel := context.WithCancel(context.TODO())
	cancel()
	return ctx
}

type causer struct{ error }

func (causer causer) Cause() error { return causer.error }

type expected struct {
	attempts uint
	error    error
}

type key struct{}

type layer struct{ error }

func (layer layer) Unwrap() error { return layer.error }

type testCase struct {
	name       string
	breaker    Breaker
	strategies How
	action     func(context.Context) error
	expected   expected
}

var testCases = []testCase{
	{
		"successful action call",
		breaker(),
		How{strategy.Wait(time.Hour)},
		func(context.Context) error { return nil },
		expected{1, nil},
	},
	{
		"failed action call",
		breaker(),
		How{strategy.Limit(10)},
		func(context.Context) error { return layer{causer{Error("failure")}} },
		expected{10, layer{causer{Error("failure")}}},
	},
	{
		"action call with interrupted breaker",
		interrupted(),
		How{strategy.Delay(time.Hour)},
		func(context.Context) error { return Error("zero iterations") },
		expected{0, context.Canceled},
	},
	{
		"have no action call",
		breaker(),
		How{strategy.Limit(0)},
		func(context.Context) error { return layer{causer{Error("failure")}} },
		expected{0, Error("have no any try")},
	},
}
