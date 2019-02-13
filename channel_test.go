package retry_test

import (
	"os"
	"testing"
	"time"

	. "github.com/kamilsk/retry/v3"
)

func TestMultiplex(t *testing.T) {
	sleep := 100 * time.Millisecond

	start := time.Now()
	<-Multiplex(WithSignal(os.Interrupt), WithTimeout(sleep))
	end := time.Now()

	if expected, obtained := sleep, end.Sub(start); expected > obtained {
		t.Errorf("an unexpected sleep time. expected: %v; obtained: %v", expected, obtained)
	}
}

func TestMultiplex_WithoutChannels(t *testing.T) {
	<-Multiplex()
}

func TestWithDeadline(t *testing.T) {
	tests := []struct {
		name     string
		deadline time.Duration
	}{
		{"normal case", 10 * time.Millisecond},
		{"past deadline", -time.Nanosecond},
	}
	for _, test := range tests {
		start := time.Now()
		<-WithDeadline(start.Add(test.deadline))
		end := time.Now()

		if !end.After(start.Add(test.deadline)) {
			t.Errorf("%s: an unexpected deadline", test.name)
		}
	}
}

func TestWithSignal_NilSignal(t *testing.T) {
	<-WithSignal(nil)
}

func TestWithTimeout(t *testing.T) {
	sleep := 100 * time.Millisecond

	start := time.Now()
	<-WithTimeout(sleep)
	end := time.Now()

	if expected, obtained := sleep, end.Sub(start); expected > obtained {
		t.Errorf("an unexpected sleep time. expected: %v; obtained: %v", expected, obtained)
	}
}
