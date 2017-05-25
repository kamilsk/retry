package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

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
	// NoColor deprecates colorize logger' output.
	// Can be changed by `-ldflags "-X 'main.NoColor=...'"`.
	NoColor = false
	// Version will always be the name of the current Git tag.
	Version string
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error occurred %q \n", r)
			os.Exit(1)
		}
	}()

	var stderrFlag, stdoutFlag int
	if Debug {
		stderrFlag = log.Lshortfile
	}

	l := &logger{
		stderr:  log.New(os.Stderr, "", stderrFlag),
		stdout:  log.New(os.Stdout, "", stdoutFlag),
		debug:   Debug,
		colored: !NoColor,
	}

	var (
		start   time.Time
		started bool
	)

	done := make(chan struct{})
	timeout, args, strategies := parse(os.Args[1:]...)
	action := func(attempt uint) error {
		if !started {
			start = time.Now()
			started = true
		} else {
			l.Infof("#%d attempt at %s... \n", attempt+1, time.Now().Sub(start))
		}
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, &buf{c: done, w: os.Stderr}
		return cmd.Run()
	}
	ctx, cancel := ctx(timeout)
	if err := retry.Retry(ctx, action, strategies...); err != nil {
		l.Errorf("error occurred: %q", err)
		close(done)
	}
	cancel()
}
