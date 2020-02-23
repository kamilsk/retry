package strategy

import (
	"errors"
	"net"
	"testing"
)

func TestCheckNetworkError(t *testing.T) {
	tests := map[string]struct {
		error    error
		expected bool
	}{
		"nil error": {
			nil,
			true,
		},
		"network address error": {
			&net.AddrError{},
			false,
		},
		"unknown network error": {
			net.UnknownNetworkError("test"),
			false,
		},
		"temporary dns error": {
			&net.DNSError{IsTemporary: true},
			true,
		},
		"not network error": {
			errors.New("test"),
			true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.expected != CheckNetworkError(1, test.error) {
				t.Errorf("strategy expected to return %v", test.expected)
			}
		})
	}
}
