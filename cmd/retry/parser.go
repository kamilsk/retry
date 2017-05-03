package main

import (
	"errors"
	"flag"
	"regexp"
	"time"

	cflag "github.com/kamilsk/retry/cmd/retry/flag" // custom flag

	"github.com/kamilsk/retry/backoff"
	"github.com/kamilsk/retry/jitter"
	"github.com/kamilsk/retry/strategy"
)

var (
	re = regexp.MustCompile(`^(\w+)(?:\[((?:[\w\.]+,?)+)\])?$`)

	compliance map[string]struct {
		cursor  interface{}
		usage   string
		handler func(*flag.Flag) (strategy.Strategy, error)
	}
	algorithms map[string]func(args string) (backoff.Algorithm, error)
	transforms map[string]func(args string) (jitter.Transformation, error)
	usage      func()
)

func parse(arguments []string) (time.Duration, []string, []strategy.Strategy) {
	cl := cflag.NewSet("retry")
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
	if err := cl.Parse(arguments); err != nil {
		panic(err)
	}

	timeout, err := time.ParseDuration(Timeout)
	if err != nil {
		panic(err)
	}

	strategies, err := handle(cl.Sequence())
	if err != nil {
		panic(err)
	}

	args := cl.Args()
	if len(args) == 0 {
		panic("please provide a command to retry")
	}

	return timeout, cl.Args(), strategies
}

func handle(flags []*flag.Flag) ([]strategy.Strategy, error) {
	strategies := make([]strategy.Strategy, 0, len(flags))

	for _, f := range flags {
		if c, ok := compliance[f.Name]; ok {
			s, err := c.handler(f)
			if err != nil {
				return nil, err
			}
			strategies = append(strategies, s)
		}
	}

	return strategies, nil
}

func parseAlgorithm(args string) (backoff.Algorithm, error) {
	m := re.FindStringSubmatch(args)
	if len(m) < 2 {
		return nil, errors.New("invalid argument " + args)
	}
	algorithm, ok := algorithms[m[1]]
	if !ok {
		return nil, errors.New("unknown algorithm " + m[1])
	}
	args = ""
	if len(m) == 3 {
		args = m[2]
	}
	return algorithm(args)
}

func parseTransform(args string) (jitter.Transformation, error) {
	m := re.FindStringSubmatch(args)
	if len(m) < 2 {
		return nil, errors.New("invalid argument " + args)
	}
	transformation, ok := transforms[m[1]]
	if !ok {
		return nil, errors.New("unknown transformation " + m[1])
	}
	args = ""
	if len(m) == 3 {
		args = m[2]
	}
	return transformation(args)
}
