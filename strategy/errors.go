package strategy

import "net"

// CheckNetworkError returns true if the error is the temporary network error.
// It also returns true if the error is not the network error.
func CheckNetworkError(_ uint, err error) bool {
	if err == nil {
		return true
	}
	if err, is := err.(net.Error); is {
		return err.Temporary() || err.Timeout()
	}
	return true
}
