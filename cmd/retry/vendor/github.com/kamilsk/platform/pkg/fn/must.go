package fn

import "github.com/pkg/errors"

// Must execs the actions step by step and raises a panic
// with error and its stack trace if something went wrong.
//
//  func Configure(cmd *cobra.Command) {
//  	Must(func() error { return cmd.MarkFlagRequired("format") })
//  }
//
func Must(actions ...func() error) {
	for _, action := range actions {
		if err := errors.WithStack(action()); err != nil {
			panic(err)
		}
	}
}
