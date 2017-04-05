// +build !go1.7

package main

import (
	"os"
	"time"

	pkg_backoff "github.com/kamilsk/retrier/backoff"
	flag2 "github.com/kamilsk/retrier/cmd/retry/flag"
	pkg_strategy "github.com/kamilsk/retrier/strategy"
	"golang.org/x/net/context"
)

func parse() (context.Context, []string, []pkg_strategy.Strategy) {
	cl := flag2.NewFlagSet("retry")
	for name, cfg := range compliance {
		cl.StringVar(cfg.cursor, name, "", cfg.usage)
	}
	cl.Parse(os.Args[1:])

	return context.Background(), cl.Args(), []pkg_strategy.Strategy{
		pkg_strategy.Limit(3),
		pkg_strategy.Backoff(pkg_backoff.Linear(10 * time.Millisecond)),
	}
}
