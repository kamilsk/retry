package retry

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	t.Run("with cancel", func(t *testing.T) {
		t.Parallel()

		var (
			sig = make(chan struct{})
			ctx = context.Context(lite{context.TODO(), sig})
		)
		if ctx.Err() != nil {
			t.Error("invalid state")
		}

		ctx, cancel := context.WithCancel(ctx)
		if ctx.Err() != nil {
			t.Error("invalid state")
		}

		verify(t, ctx, cancel, sig)
	})

	t.Run("with deadline", func(t *testing.T) {
		t.Parallel()

		var (
			sig = make(chan struct{})
			ctx = context.Context(lite{context.TODO(), sig})
		)
		if ctx.Err() != nil {
			t.Error("invalid state")
		}

		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Hour))
		if ctx.Err() != nil {
			t.Error("invalid state")
		}

		verify(t, ctx, cancel, sig)
	})

	t.Run("with timeout", func(t *testing.T) {
		t.Parallel()

		var (
			sig = make(chan struct{})
			ctx = context.Context(lite{context.TODO(), sig})
		)
		if ctx.Err() != nil {
			t.Error("invalid state")
		}

		ctx, cancel := context.WithTimeout(ctx, time.Hour)
		if ctx.Err() != nil {
			t.Error("invalid state")
		}

		verify(t, ctx, cancel, sig)
	})

	t.Run("with value", func(t *testing.T) {
		t.Parallel()

		var (
			sig = make(chan struct{})
			ctx = context.Context(lite{context.TODO(), sig})
		)
		if ctx.Err() != nil {
			t.Error("invalid state")
		}

		ctx = context.WithValue(ctx, key{}, "value")
		if expected, obtained := "value", ctx.Value(key{}); obtained != expected {
			t.Errorf("expected: %q, obtained: %q", expected, obtained)
		}

		close(sig)
	})
}

// helpers

func stop(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
}

func verify(t *testing.T, ctx context.Context, cancel context.CancelFunc, sig chan struct{}) {
	t.Helper()

	timer := time.NewTimer(schedule)
	close(sig)
	select {
	case <-timer.C:
		t.Error("invalid state")
	case <-ctx.Done():
		if ctx.Err() == nil {
			t.Error("invalid state")
		}
	}

	stop(timer)
	cancel()
}

type key struct{}

var schedule = 10 * time.Duration(runtime.NumCPU()) * time.Millisecond
