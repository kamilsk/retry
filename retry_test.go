package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/kamilsk/retry/v4"
	. "github.com/kamilsk/retry/v4/strategy"
	"github.com/stretchr/testify/assert"
)

var delta = 10 * time.Millisecond

func TestRetry(t *testing.T) {
	type Assert struct {
		Attempts uint
		Error    func(assert.TestingT, error, ...interface{}) bool
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
			[]func(attempt uint, err error) bool{Delay(delta), Limit(10000)},
			errors.New("zero iterations"),
			Assert{0, func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.Error(t, err) && assert.True(t, IsInterrupted(err))
			}},
		},
		{
			"one iteration",
			nil,
			nil,
			nil,
			Assert{1, assert.NoError},
		},
		{
			"two iterations",
			nil,
			[]func(attempt uint, err error) bool{Limit(2)},
			errors.New("two iterations"),
			Assert{2, assert.Error},
		},
		{
			"three iterations",
			newBreaker(),
			[]func(attempt uint, err error) bool{Limit(3)},
			errors.New("three iterations"),
			Assert{3, func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.EqualError(t, err, "three iterations")
			}},
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
			tc.assert.Error(t, err)
			_, is := IsRecovered(err)
			assert.False(t, is)
			assert.Equal(t, tc.assert.Attempts, total)
		})
	}
	t.Run("unexpected panic", func(t *testing.T) {
		err := Retry(newBreaker(), func(uint) error { panic("Catch Me If You Can") })
		assert.Error(t, err)
		cause, is := IsRecovered(err)
		assert.True(t, is)
		assert.Equal(t, "Catch Me If You Can", cause)
	})
}

func TestTry(t *testing.T) {
	type Assert struct {
		Attempts uint
		Error    func(assert.TestingT, error, ...interface{}) bool
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
			[]func(attempt uint, err error) bool{Delay(delta), Limit(10000)},
			errors.New("zero iterations"),
			Assert{0, func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.Error(t, err) && assert.True(t, IsInterrupted(err))
			}},
		},
		{
			"one iteration",
			nil,
			nil,
			nil,
			Assert{1, assert.NoError},
		},
		{
			"two iterations",
			nil,
			[]func(attempt uint, err error) bool{Limit(2)},
			errors.New("two iterations"),
			Assert{2, assert.Error},
		},
		{
			"three iterations",
			newPanicBreaker(),
			[]func(attempt uint, err error) bool{Limit(3)},
			errors.New("three iterations"),
			Assert{3, func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.EqualError(t, err, "three iterations")
			}},
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
			tc.assert.Error(t, err)
			_, is := IsRecovered(err)
			assert.False(t, is)
			assert.Equal(t, tc.assert.Attempts, total)
		})
	}
	t.Run("unexpected panic", func(t *testing.T) {
		err := Try(newBreaker(), func(uint) error { panic("Catch Me If You Can") })
		assert.Error(t, err)
		cause, is := IsRecovered(err)
		assert.True(t, is)
		assert.Equal(t, "Catch Me If You Can", cause)
	})
}

func TestTryContext(t *testing.T) {
	type Assert struct {
		Attempts uint
		Error    func(assert.TestingT, error, ...interface{}) bool
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
			[]func(attempt uint, err error) bool{Delay(delta), Limit(10000)},
			errors.New("zero iterations"),
			Assert{0, func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.Error(t, err) && assert.True(t, IsInterrupted(err))
			}},
		},
		{
			"one iteration",
			nil,
			nil,
			nil,
			Assert{1, assert.NoError},
		},
		{
			"two iterations",
			nil,
			[]func(attempt uint, err error) bool{Limit(2)},
			errors.New("two iterations"),
			Assert{2, assert.Error},
		},
		{
			"three iterations",
			context.Background(),
			[]func(attempt uint, err error) bool{Limit(3)},
			errors.New("three iterations"),
			Assert{3, func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.EqualError(t, err, "three iterations")
			}},
		},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			var total uint
			action := func(ctx context.Context, attempt uint) error {
				assert.Equal(t, tc.ctx, ctx)
				total = attempt + 1
				return tc.error
			}
			err := TryContext(tc.ctx, action, tc.strategies...)
			tc.assert.Error(t, err)
			_, is := IsRecovered(err)
			assert.False(t, is)
			assert.Equal(t, tc.assert.Attempts, total)
		})
	}
	t.Run("unexpected panic", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		err := TryContext(ctx, func(context.Context, uint) error { panic("Catch Me If You Can") })
		assert.Error(t, err)
		cause, is := IsRecovered(err)
		assert.True(t, is)
		assert.Equal(t, "Catch Me If You Can", cause)
		cancel()
	})
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

func newPanicBreaker() BreakCloser {
	return &panicBreaker{newBreaker()}
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

type panicBreaker struct {
	*contextBreaker
}

func (*panicBreaker) Close() {
	panic("unexpected method call")
}
