package main

import (
	"bytes"
	"testing"
)

func TestBuffer_Write(t *testing.T) {
	w := bytes.NewBuffer(nil)
	b := buf{c: make(chan struct{}), w: w}

	if _, err := b.Write([]byte("test")); err != nil {
		t.Errorf("unexpected error %q", err)
	}
	close(b.c)
	if _, err := b.Write([]byte("buffer")); err != nil {
		t.Errorf("unexpected error %q", err)
	}

	expected, obtained := "test", w.String()
	if obtained != expected {
		t.Errorf("unexpected buffer value, expected %q, obtained %q", expected, obtained)
	}
}
