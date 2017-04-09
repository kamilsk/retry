// +build go1.7

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kamilsk/retry/cmd/retry/flag"
	"github.com/kamilsk/retry/strategy"
)

func parse() (context.Context, []string, []strategy.Strategy) {
	cl := flag.NewFlagSet("retry")

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error occured %q \n", r)
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
