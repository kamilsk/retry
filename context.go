package retry

import "context"

type lite struct {
	context.Context
	signal <-chan struct{}
}

func (ctx lite) Done() <-chan struct{} {
	return ctx.signal
}

func (ctx lite) Err() error {
	select {
	case <-ctx.signal:
		return context.Canceled
	default:
		return nil
	}
}

// equal to go.octolab.org/errors.Unwrap
func unwrap(err error) error {
	// compatible with github.com/pkg/errors
	type causer interface {
		Cause() error
	}
	// compatible with built-in errors since 1.13
	type wrapper interface {
		Unwrap() error
	}

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
