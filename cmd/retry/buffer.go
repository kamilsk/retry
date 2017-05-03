package main

import (
	"io"
	"io/ioutil"
)

type buf struct {
	c chan struct{}
	w io.Writer
}

func (b *buf) Write(p []byte) (n int, err error) {
	select {
	case <-b.c:
		return ioutil.Discard.Write(p)
	default:
		return b.w.Write(p)
	}
}
