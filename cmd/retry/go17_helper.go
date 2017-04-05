// +build go1.7

package main

import (
	"context"
	"os"
	"time"

	pkg_backoff "github.com/kamilsk/retrier/backoff"
	flag2 "github.com/kamilsk/retrier/cmd/retry/flag"
	pkg_strategy "github.com/kamilsk/retrier/strategy"
)

func parse() (context.Context, []string, []pkg_strategy.Strategy) {
	cl := flag2.NewFlagSet("retry")
	for name, cfg := range compliance {
		switch cursor := cfg.cursor.(type) {
		case *string:
			cl.StringVar(cursor, name, "", cfg.usage)
		case *bool:
			cl.BoolVar(cursor, name, false, cfg.usage)
		}

	}
	cl.Parse(os.Args[1:])

	return context.Background(), cl.Args(), []pkg_strategy.Strategy{
		pkg_strategy.Limit(3),
		pkg_strategy.Backoff(pkg_backoff.Linear(10 * time.Millisecond)),
	}
}
