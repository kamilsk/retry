package retry_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/strategy"
)

func TestDo(t *testing.T) {
	type expected struct {
		attempts uint
		error    error
	}

	tests := map[string]struct {
		breaker    strategy.Breaker
		strategies How
		action     func(context.Context) error
		expected   expected
	}{
		"successful action call": {
			breaker(),
			How{strategy.Wait(time.Hour)},
			func(context.Context) error { return nil },
			expected{1, nil},
		},
		"failed action call": {
			breaker(),
			How{strategy.Limit(10)},
			func(context.Context) error { return layer{causer{errors.New("failure")}} },
			expected{10, layer{causer{errors.New("failure")}}},
		},
		"action call with interrupted breaker": {
			interrupted(),
			How{strategy.Delay(time.Hour)},
			func(context.Context) error { return errors.New("zero iterations") },
			expected{0, context.Canceled},
		},
		"have no action call": {
			breaker(),
			How{strategy.Limit(0)},
			func(context.Context) error { return layer{causer{errors.New("failure")}} },
			expected{0, Error("have no any try")},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
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
		ctx := context.WithValue(context.TODO(), key{}, "value")
		action := func(ctx context.Context) error {
			if !reflect.DeepEqual("value", ctx.Value(key{})) {
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
	type expected struct {
		attempts uint
		error    error
	}

	tests := map[string]struct {
		breaker    strategy.Breaker
		strategies How
		action     func(context.Context) error
		expected   expected
	}{
		"successful action call": {
			breaker(),
			How{strategy.Wait(time.Hour)},
			func(context.Context) error { return nil },
			expected{1, nil},
		},
		"failed action call": {
			breaker(),
			How{strategy.Limit(10)},
			func(context.Context) error { return layer{causer{errors.New("failure")}} },
			expected{10, layer{causer{errors.New("failure")}}},
		},
		"action call with interrupted breaker": {
			interrupted(),
			How{strategy.Delay(time.Hour)},
			func(context.Context) error { return errors.New("zero iterations") },
			expected{0, context.Canceled},
		},
		"have no action call": {
			breaker(),
			How{strategy.Limit(0)},
			func(context.Context) error { return layer{causer{errors.New("failure")}} },
			expected{0, Error("have no any try")},
		},
		"action call with error panic": {
			breaker(),
			How{strategy.Wait(time.Hour)},
			func(context.Context) error { panic(errors.New("failure")) },
			expected{1, errors.New("failure")},
		},
		"action call with non-error panic": {
			breaker(),
			How{strategy.Wait(time.Hour)},
			func(context.Context) error { panic("non-error") },
			expected{1, fmt.Errorf("retry: unexpected panic: %#v", "non-error")},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
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
		ctx := context.WithValue(context.TODO(), key{}, "value")
		action := func(ctx context.Context) error {
			if !reflect.DeepEqual("value", ctx.Value(key{})) {
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

type key struct{}

func breaker() strategy.Breaker {
	return context.TODO()
}

func interrupted() strategy.Breaker {
	ctx, cancel := context.WithCancel(context.TODO())
	cancel()
	return ctx
}

type layer struct{ error }

func (layer layer) Unwrap() error { return layer.error }

type causer struct{ error }

func (causer causer) Cause() error { return causer.error }
