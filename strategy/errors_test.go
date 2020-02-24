package strategy

import (
	"errors"
	"math/rand"
	"net"
	"testing"
)

func TestCheckError(t *testing.T) {
	generator := rand.New(rand.NewSource(0))

	tests := map[string]struct {
		error    error
		defaults bool
		expected bool
	}{
		"nil error": {
			nil,
			Skip,
			true,
		},
		"retriable error": {
			exampleError("test"),
			Strict,
			true,
		},
		"not retriable error": {
			errors.New("test"),
			Skip,
			true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			policy := CheckError(Skip)
			if test.expected != policy(uint(generator.Uint32()), test.error) {
				t.Errorf("strategy expected to return %v", test.expected)
			}
		})
	}
}

func TestCheckNetworkError(t *testing.T) {
	generator := rand.New(rand.NewSource(0))

	tests := map[string]struct {
		error    error
		defaults bool
		expected bool
	}{
		"nil error": {
			nil,
			Skip,
			true,
		},
		"network address error": {
			&net.AddrError{},
			Strict,
			false,
		},
		"unknown network error": {
			net.UnknownNetworkError("test"),
			Strict,
			false,
		},
		"temporary dns error": {
			&net.DNSError{IsTemporary: true},
			Strict,
			true,
		},
		"not network error": {
			errors.New("test"),
			Skip,
			true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			policy := CheckNetworkError(Skip)
			if test.expected != policy(uint(generator.Uint32()), test.error) {
				t.Errorf("strategy expected to return %v", test.expected)
			}
		})
	}
}

type exampleError string

func (err exampleError) Error() string   { return string(err) }
func (err exampleError) Retriable() bool { return true }
