// +build !go1.7

package main

import (
	"time"

	"golang.org/x/net/context"
)

func ctx(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
