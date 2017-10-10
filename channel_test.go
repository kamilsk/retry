package retry_test

import (
	"os"
	"testing"
	"time"

	"github.com/kamilsk/retry"
)

func TestMultiplex(t *testing.T) {
	sleep := 100 * time.Millisecond

	start := time.Now()
	<-retry.Multiplex(retry.WithSignal(os.Interrupt), retry.WithTimeout(sleep))
	end := time.Now()

	if expected, obtained := sleep, end.Sub(start); expected > obtained {
		t.Errorf("an unexpected sleep time. expected: %v; obtained: %v", expected, obtained)
	}
}

func TestMultiplex_WithoutChannels(t *testing.T) {
	<-retry.Multiplex()
}

func TestWithDeadline(t *testing.T) {
	sleep := time.Now().Add(100 * time.Millisecond)

	<-retry.WithDeadline(sleep)
	end := time.Now()

	if expected, obtained := sleep, end; expected.After(obtained) {
		t.Errorf("an unexpected sleep time. expected: %v; obtained: %v", expected, obtained)
	}
}

func TestWithSignal_NilSignal(t *testing.T) {
	<-retry.WithSignal(nil)
}

func TestWithTimeout(t *testing.T) {
	sleep := 100 * time.Millisecond

	start := time.Now()
	<-retry.WithTimeout(sleep)
	end := time.Now()

	if expected, obtained := sleep, end.Sub(start); expected > obtained {
		t.Errorf("an unexpected sleep time. expected: %v; obtained: %v", expected, obtained)
	}
}
