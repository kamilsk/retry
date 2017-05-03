// +build go1.7

package main

import (
	"context"
	"time"
)

func ctx(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
