package sandbox

import "net"

const (
	Skip = true
	Stop = false
)

// A Breaker carries a cancellation signal to interrupt an action execution.
//
// It is a subset of the built-in context and github.com/kamilsk/breaker interfaces.
type Breaker = interface {
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}
	// If Done is not yet closed, Err returns nil.
	// If Done is closed, Err returns a non-nil error.
	// After Err returns a non-nil error, successive calls to Err return the same error.
	Err() error
}

// ErrorHandler defines a function that CheckError calls
// to determine whether it should make the next attempt or not.
// Returning true allows for the next attempt to be made.
// Returning false halts the retrying process and returns the last error
// returned by the called Action.
type ErrorHandler = func(error) bool

// CheckError creates a Strategy that checks an error and returns
// if an error is retriable or not. Otherwise, it returns the defaults.
func CheckError(handlers ...func(error) bool) func(Breaker, uint, error) bool {
	// equal to go.octolab.org/errors.Retriable
	type retriable interface {
		error
		Retriable() bool // Is the error retriable?
	}

	return func(_ Breaker, _ uint, err error) bool {
		if err == nil {
			return true
		}
		if err, is := err.(retriable); is {
			return err.Retriable()
		}
		for _, handle := range handlers {
			if !handle(err) {
				return false
			}
		}
		return true
	}
}

// NetworkError creates an error Handler that checks an error and returns true
// if an error is the temporary network error.
// The Handler returns the defaults if an error is not a network error.
func NetworkError(defaults bool) func(error) bool {
	return func(err error) bool {
		if err, is := err.(net.Error); is {
			return err.Temporary() || err.Timeout()
		}
		return defaults
	}
}
