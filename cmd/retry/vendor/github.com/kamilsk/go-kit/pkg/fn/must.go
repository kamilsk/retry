package fn

import "github.com/pkg/errors"

// Must execs actions step by step and raises a panic
// with error and its stack trace if something went wrong.
func Must(actions ...func() error) {
	for _, action := range actions {
		if err := errors.WithStack(action()); err != nil {
			panic(err)
		}
	}
}
