package net

import (
	"net"

	"github.com/kamilsk/retry/strategy"
)

// CheckNetError creates a Strategy that will check if network request failed with a temporary error or timing.
func CheckNetError() strategy.ExtendedStrategy {
	return func(attempt uint, err error) bool {
		if err, ok := err.(net.Error); ok {
			return err.Timeout() || err.Temporary()
		}
		return true
	}
}
