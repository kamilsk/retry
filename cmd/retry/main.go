package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/kamilsk/retry"
	"github.com/kamilsk/retry/cmd/retry/flag"
	"github.com/kamilsk/retry/strategy"
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
	// Can be changed by `-ldflags "-X 'main.Timeout=...'"`.
	NoColor = false
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
	timeout, args, strategies := parse(os.Args[1:]...)
	action := func(attempt uint) error {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, &buf{c: done, w: os.Stderr}
		return cmd.Run()
	}
	ctx, cancel := ctx(timeout)
	if err := retry.Retry(ctx, action, strategies...); err != nil {
		l.Errorf("error occcured: %q", err)
		close(done)
	}
	cancel()
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
