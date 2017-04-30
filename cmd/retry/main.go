package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/kamilsk/retry"
)

var (
	// Debug prints verbose information to stdout.
	// Can be changed by `-ldflags "-X" 'main.Debug=..."'`
	// or `-v` parameter.
	Debug = false
	// Timeout is a timeout of retried operation.
	// Can be changed by `-ldflags "-X 'main.Timeout=...'"`
	// or `-timeout ...` parameter.
	Timeout = "1m"
	// Version will always be the name of the current Git tag.
	Version string
)

var l *logger

func init() {
	l = &logger{
		stderr: log.New(os.Stderr, "", log.Lshortfile),
		stdout: log.New(os.Stdout, "", 0),
		debug:  Debug,
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error occurred %q \n", r)
			os.Exit(1)
		}
	}()

	done := make(chan struct{})
	timeout, args, strategies := parse()
	action := func(attempt uint) error {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, &buf{c: done, w: os.Stderr}
		return cmd.Run()
	}
	ctx, cancel := ctx(timeout)
	if err := retry.Retry(ctx, action, strategies...); err != nil {
		l.Errorf("error occcured %q", err)
		close(done)
	}
	cancel()
}
