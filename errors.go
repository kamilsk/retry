package retry

const internal Error = "have no any try"

// Error defines a string-based error without a different root cause.
type Error string

// Error returns a string representation of an error.
func (err Error) Error() string { return string(err) }

// Unwrap always returns nil means that an error doesn't have other root cause.
func (err Error) Unwrap() error { return nil }

func unwrap(err error) error {
	for err != nil {
		layer, is := err.(wrapper)
		if is {
			err = layer.Unwrap()
			continue
		}
		cause, is := err.(causer)
		if is {
			err = cause.Cause()
			continue
		}
		break
	}
	return err
}

// compatible with github.com/pkg/errors
type causer interface {
	Cause() error
}

// compatible with built-in errors since 1.13
type wrapper interface {
	Unwrap() error
}
