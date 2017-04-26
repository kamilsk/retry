package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/kamilsk/retry"
)

var (
	// Timeout is a timeout of retried operation.
	// Can be changed by `-ldflags "-X 'main.Timeout=...'"` or `-timeout ...` parameter.
	Timeout = "1m"
	// Version will always be the name of the current Git tag.
	Version string
)

func main() {
	done := make(chan struct{})
	ctx, cancel, args, strategies := parse()
	defer cancel()
	action := func(attempt uint) error {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, &buf{c: done, w: os.Stderr}
		return cmd.Run()
	}
	if err := retry.Retry(ctx, action, strategies...); err != nil {
		fmt.Fprintf(os.Stderr, "error occurred %q \n", err)
		close(done)
	}
}

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
