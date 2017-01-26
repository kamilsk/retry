package net

import (
	"net"

	"github.com/kamilsk/retrier/strategy"
)

// CheckNetworkError creates a Strategy that will check if network request failed with a temporary error or timing.
func CheckNetworkError() strategy.Strategy {
	return func(attempt uint, err error) bool {
		if err, ok := err.(net.Error); ok {
			return err.Timeout() || err.Temporary()
		}
		return true
	}
}
