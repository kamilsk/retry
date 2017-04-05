package main

import (
	"flag"

	"github.com/kamilsk/retrier/strategy"
)

type Compliance map[string]struct {
	cursor  interface{}
	usage   string
	handler Handler
}

type Handler interface {
	Handle(flag.Flag) strategy.Strategy
}

var compliance Compliance

// TODO: generate it

func init() {
	var (
		infinite                        bool
		limit, delay, backoff, tbackoff string
	)
	compliance = Compliance{
		"infinite": {cursor: &infinite, usage: ""},
		"limit":    {cursor: &limit, usage: ""},
		"delay":    {cursor: &delay, usage: ""},
		"backoff":  {cursor: &backoff, usage: ""},
		"tbackoff": {cursor: &tbackoff, usage: ""},
	}
}
