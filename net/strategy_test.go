package net

import (
	"net"
	"testing"
)

func TestCheckNetError(t *testing.T) {
	strategy := CheckNetError()

	if !strategy(0, nil) {
		t.Error("strategy expected to return true")
	}

	if !strategy(0, &net.DNSError{IsTimeout: true}) {
		t.Error("strategy expected to return true")
	}

	if strategy(0, &net.DNSError{}) {
		t.Error("strategy expected to return false")
	}
}
