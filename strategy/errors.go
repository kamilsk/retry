package strategy

import "net"

const (
	Skip   = true
	Strict = false
)

// CheckNetworkError creates a Strategy that checks an error and returns true
// if an error is the temporary network error.
// The Strategy returns the defaults if an error is not a network error.
func CheckNetworkError(defaults bool) Strategy {
	return func(_ uint, err error) bool {
		if err == nil {
			return true
		}
		if err, is := err.(net.Error); is {
			return err.Temporary() || err.Timeout()
		}
		return defaults
	}
}
