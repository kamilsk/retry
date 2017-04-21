package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/kamilsk/retry"
	"github.com/kamilsk/retry/cmd/retry/flag"
	"github.com/kamilsk/retry/strategy"
)

// Timeout is a timeout of retried operation.
// Can be changed by `-ldflags "-X 'main.Timeout=...'"` or `-timeout ...` parameter.
var Timeout = "1m"

func main() {
	done := make(chan struct{})
	ctx, args, strategies := parse()
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

func parse() (context.Context, []string, []strategy.Strategy) {
	cl := flag.NewFlagSet("retry")

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error occurred %q \n", r)
			cl.Usage()
			os.Exit(1)
		}
	}()

	cl.Usage = usage
	for name, cfg := range compliance {
		switch cursor := cfg.cursor.(type) {
		case *string:
			cl.StringVar(cursor, name, "", cfg.usage)
		case *bool:
			cl.BoolVar(cursor, name, false, cfg.usage)
		}

	}
	cl.StringVar(&Timeout, "timeout", Timeout, "value which supported by time.ParseDuration")
	cl.Parse(os.Args[1:])

	timeout, err := time.ParseDuration(Timeout)
	if err != nil {
		panic(err)
	}

	strategies, err := handle(cl.Flags())
	if err != nil {
		panic(err)
	}

	args := cl.Args()
	if len(args) == 0 {
		panic("please provide a command to retry")
	}

	ctx, _ := context.WithTimeout(context.Background(), timeout)

	return ctx, cl.Args(), strategies
}
