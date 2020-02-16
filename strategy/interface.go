package strategy

// A Breaker carries a cancellation signal to break an action execution.
//
// It is a subset of context.Context and github.com/kamilsk/breaker.Breaker.
type Breaker interface {
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}
	// If Done is not yet closed, Err returns nil.
	// If Done is closed, Err returns a non-nil error
	// related to an occurred cancellation signal.
	// After Err returns a non-nil error, successive calls to Err
	// return the same error.
	Err() error
}
