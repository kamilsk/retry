// TODO: generate it

package main

import (
	"context"
	"flag"
	"reflect"

	"time"

	pkg_backoff "github.com/kamilsk/retrier/backoff"
	pkg_jitter "github.com/kamilsk/retrier/jitter"
	pkg_strategy "github.com/kamilsk/retrier/strategy"
)

var (
	strategies map[string]reflect.Value
	algorithms map[string]reflect.Value
	transforms map[string]reflect.Value
)

func init() {
	strategies = map[string]reflect.Value{
		"infinite": reflect.ValueOf(pkg_strategy.Infinite),
		"limit":    reflect.ValueOf(pkg_strategy.Limit),
		"delay":    reflect.ValueOf(pkg_strategy.Delay),
		"backoff":  reflect.ValueOf(pkg_strategy.Backoff),
		"tbackoff": reflect.ValueOf(pkg_strategy.BackoffWithJitter),
	}
	algorithms = map[string]reflect.Value{
		"inc":    reflect.ValueOf(pkg_backoff.Incremental),
		"lin":    reflect.ValueOf(pkg_backoff.Linear),
		"epx":    reflect.ValueOf(pkg_backoff.Exponential),
		"binexp": reflect.ValueOf(pkg_backoff.BinaryExponential),
		"fib":    reflect.ValueOf(pkg_backoff.Fibonacci),
	}
	transforms = map[string]reflect.Value{
		"full":  reflect.ValueOf(pkg_jitter.Full),
		"equal": reflect.ValueOf(pkg_jitter.Equal),
		"dev":   reflect.ValueOf(pkg_jitter.Deviation),
		"ndist": reflect.ValueOf(pkg_jitter.NormalDistribution),
	}
}

func parse() (context.Context, []string, []pkg_strategy.Strategy) {
	var infinite, limit, delay, backoff, tbackoff string

	flag.CommandLine.Init("retry", flag.PanicOnError)
	flag.StringVar(&infinite, "infinite", "", "")
	flag.StringVar(&limit, "limit", "", "")
	flag.StringVar(&delay, "delay", "", "")
	flag.StringVar(&backoff, "backoff", "", "")
	flag.StringVar(&tbackoff, "tbackoff", "", "")
	flag.Parse()

	return context.Background(), flag.Args(), []pkg_strategy.Strategy{
		pkg_strategy.Limit(3),
		pkg_strategy.Backoff(pkg_backoff.Linear(10 * time.Millisecond)),
	}
}
