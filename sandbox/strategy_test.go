package sandbox_test

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"testing"

	. "github.com/kamilsk/retry/sandbox"
)

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
			[]func(error) bool{NetworkError(Stop)},
			true,
		},
		"retriable error without handlers": {
			retriable("yes"),
			nil,
			true,
		},
		"retriable error with strict handlers": {
			retriable("yes"),
			[]func(error) bool{NetworkError(Stop)},
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
			[]func(error) bool{NetworkError(Stop)},
			false,
		},
		"network address error without strict check": {
			&net.AddrError{},
			[]func(error) bool{NetworkError(Skip)},
			false,
		},
		"temporary dns error": {
			&net.DNSError{IsTemporary: true},
			[]func(error) bool{NetworkError(Stop)},
			true,
		},
		"an error with strict check": {
			errors.New("test"),
			[]func(error) bool{NetworkError(Stop)},
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
			if obtained := policy(context.TODO(), attempt, test.error); test.expected != obtained {
				t.Errorf("expected: %v, obtained: %v", test.expected, obtained)
			}
		})
	}
}

// helpers

type retriable string

func (err retriable) Error() string   { return string(err) }
func (err retriable) Retriable() bool { return err == "yes" }
