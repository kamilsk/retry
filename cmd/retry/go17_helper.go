// +build go1.7

package main

import (
	"context"
	"os"
	"time"

	"github.com/kamilsk/retrier/backoff"
	"github.com/kamilsk/retrier/cmd/retry/flag"
	"github.com/kamilsk/retrier/strategy"
)

func parse() (context.Context, []string, []strategy.Strategy) {
	cl := flag.NewFlagSet("retry")
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
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	return ctx, cl.Args(), []strategy.Strategy{
		strategy.Limit(3),
		strategy.Backoff(backoff.Linear(10 * time.Millisecond)),
	}
}
