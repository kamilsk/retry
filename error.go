package retry

// Error defines package errors.
type Error string

// Error implements the error interface.
func (err Error) Error() string {
	return string(err)
}

const (
	Panicked    Error = "panic unexpected"
	Interrupted Error = "operation interrupted"
)

// IsInterrupted checks that the error is related to the Breaker interruption.
// Deprecated: use err == retry.Interrupted instead.
func IsInterrupted(err error) bool {
	return err == Interrupted
}

// IsRecovered checks that the error is related to unhandled Action's panic
// and returns an original cause of panic.
func IsRecovered(err error) (interface{}, bool) {
	if h, is := err.(panicHandler); is {
		return h.recovered, true
	}
	return nil, false
}

type panicHandler struct {
	error
	recovered interface{}
}

// Cause returns the underlying cause of the error.
// Friendly for the github.com/pkg/errors package.
func (h panicHandler) Cause() error {
	return h.error
}

func (panicHandler) recover(err *error) {
	if r := recover(); r != nil {
		*err = panicHandler{Panicked, r}
	}
}
