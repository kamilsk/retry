package main

import (
	"flag"

	"github.com/kamilsk/retrier/strategy"
)

type Compliance map[string]struct {
	cursor  interface{}
	usage   string
	handler func(*flag.Flag) (strategy.Strategy, error)
}

var compliance Compliance

// TODO: generate it

func init() {
	var (
		infinite                        bool
		limit, delay, backoff, tbackoff string
	)
	compliance = Compliance{
		"infinite": {cursor: &infinite},
		"limit":    {cursor: &limit},
		"delay":    {cursor: &delay},
		"backoff":  {cursor: &backoff},
		"tbackoff": {cursor: &tbackoff},
	}
}

func InfiniteStrategy_gen(f *flag.Flag) (strategy.Strategy, error) {
	return strategy.Infinite(), nil
}

func LimitStrategy_gen(f *flag.Flag) (strategy.Strategy, error) {
	return strategy.Limit(0), nil
}

func DelayStrategy_get(f *flag.Flag) (strategy.Strategy, error) {
	return strategy.Delay(0), nil
}

func WaitStrategy_get(f *flag.Flag) (strategy.Strategy, error) {
	return strategy.Wait(0, 0, 0), nil
}

func BackoffStrategy_get(f *flag.Flag) (strategy.Strategy, error) {
	return strategy.Backoff(nil), nil
}

func BackoffWithJitterStrategy_get(f *flag.Flag) (strategy.Strategy, error) {
	return strategy.BackoffWithJitter(nil, nil), nil
}
