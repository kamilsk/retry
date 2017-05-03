// +build go1.7

package main

import (
	"testing"
	"time"
)

func TestContext_New(t *testing.T) {
	context, cancel := ctx(time.Second)

	if context.Err() != nil {
		t.Errorf("unexpected error, obtainer %q", context.Err())
	}

	cancel()
	if context.Err() == nil {
		t.Error("expected error, obtained nil")
	}
}
