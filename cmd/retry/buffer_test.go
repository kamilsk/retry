package main

import (
	"bytes"
	"testing"
)

func TestBuffer_Write(t *testing.T) {
	w := bytes.NewBuffer(nil)
	b := buf{c: make(chan struct{}), w: w}

	b.Write([]byte("test"))
	close(b.c)
	b.Write([]byte("buffer"))

	expected, obtained := "test", w.String()
	if obtained != expected {
		t.Errorf("unexpected buffer value, expected %q, obtained %q", expected, obtained)
	}
}
