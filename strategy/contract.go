package strategy

// A Breaker carries a cancellation signal to break an action execution.
//
// It is a subset of the built-in Context and github.com/kamilsk/breaker.Breaker.
type Breaker interface {
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}
	// Err returns a non-nil error if Done is closed and nil otherwise.
	// After Err returns a non-nil error, successive calls to Err return the same error.
	Err() error
}
