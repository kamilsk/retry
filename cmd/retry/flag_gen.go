package main

import (
	"errors"
	"flag"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kamilsk/retry/backoff"
	"github.com/kamilsk/retry/jitter"
	"github.com/kamilsk/retry/strategy"
)

var (
	compliance map[string]struct {
		cursor  interface{}
		usage   string
		handler func(*flag.Flag) (strategy.Strategy, error)
	}
	algorithms map[string]func(args string) (backoff.Algorithm, error)
	transforms map[string]func(args string) (jitter.Transformation, error)
	re         = regexp.MustCompile(`^(\w+)(?:\[((?:\w+,?)+)\])?$`)
)

func generator() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
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

func parseAlgorithm(arg string) (backoff.Algorithm, error) {
	m := re.FindStringSubmatch(arg)
	if len(m) < 2 {
		return nil, errors.New("invalid argument " + arg)
	}
	algorithm, ok := algorithms[m[1]]
	if !ok {
		return nil, errors.New("unknown algorithn " + m[1])
	}
	args := ""
	if len(m) == 3 {
		args = m[2]
	}
	return algorithm(args)
}

func parseTransform(arg string) (jitter.Transformation, error) {
	m := re.FindStringSubmatch(arg)
	if len(m) < 2 {
		return nil, errors.New("invalid argument " + arg)
	}
	transformation, ok := transforms[m[1]]
	if !ok {
		return nil, errors.New("unknown transformation " + m[1])
	}
	args := ""
	if len(m) == 3 {
		args = m[2]
	}
	return transformation(args)
}

// TODO: generate it

// strategies

func generatedInfiniteStrategy(_ *flag.Flag) (strategy.Strategy, error) {
	return strategy.Infinite(), nil
}

func generatedLimitStrategy(f *flag.Flag) (strategy.Strategy, error) {
	attemptLimit, err := strconv.ParseUint(f.Value.String(), 10, 0)
	if err != nil {
		return nil, err
	}
	return strategy.Limit(uint(attemptLimit)), nil
}

func generatedDelayStrategy(f *flag.Flag) (strategy.Strategy, error) {
	duration, err := time.ParseDuration(f.Value.String())
	if err != nil {
		return nil, err
	}
	return strategy.Delay(duration), nil
}

func generatedWaitStrategy(f *flag.Flag) (strategy.Strategy, error) {
	args := strings.Split(f.Value.String(), ",")
	durations := make([]time.Duration, 0, len(args))
	for _, arg := range args {
		duration, err := time.ParseDuration(arg)
		if err != nil {
			return nil, err
		}
		durations = append(durations, duration)
	}
	return strategy.Wait(durations...), nil
}

func generatedBackoffStrategy(f *flag.Flag) (strategy.Strategy, error) {
	algorithm, err := parseAlgorithm(f.Value.String())
	if err != nil {
		return nil, err
	}
	return strategy.Backoff(algorithm), nil
}

func generatedBackoffWithJitterStrategy(f *flag.Flag) (strategy.Strategy, error) {
	args := strings.Split(f.Value.String(), ",")
	if len(args) != 2 {
		return nil, errors.New("invalid argument count")
	}
	algorithm, err := parseAlgorithm(args[0])
	if err != nil {
		return nil, err
	}
	transform, err := parseTransform(args[1])
	if err != nil {
		return nil, err
	}
	return strategy.BackoffWithJitter(algorithm, transform), nil
}

// algorithms

func generatedIncrementalAlgorithm(raw string) (backoff.Algorithm, error) {
	args := strings.Split(raw, ",")
	if len(args) != 2 {
		return nil, errors.New("invalid argument count")
	}
	initial, err := time.ParseDuration(args[0])
	if err != nil {
		return nil, err
	}
	increment, err := time.ParseDuration(args[0])
	if err != nil {
		return nil, err
	}
	return backoff.Incremental(initial, increment), nil
}

func generatedLinearAlgorithm(raw string) (backoff.Algorithm, error) {
	factor, err := time.ParseDuration(raw)
	if err != nil {
		return nil, err
	}
	return backoff.Linear(factor), nil
}

func generatedExponentialAlgorithm(raw string) (backoff.Algorithm, error) {
	args := strings.Split(raw, ",")
	if len(args) != 2 {
		return nil, errors.New("invalid argument count")
	}
	factor, err := time.ParseDuration(args[0])
	if err != nil {
		return nil, err
	}
	base, err := strconv.ParseFloat(args[1], 0)
	if err != nil {
		return nil, err
	}
	return backoff.Exponential(factor, base), nil
}

func generatedBinaryExponentialAlgorithm(raw string) (backoff.Algorithm, error) {
	factor, err := time.ParseDuration(raw)
	if err != nil {
		return nil, err
	}
	return backoff.BinaryExponential(factor), nil
}

func generatedFibonacciAlgorithm(raw string) (backoff.Algorithm, error) {
	factor, err := time.ParseDuration(raw)
	if err != nil {
		return nil, err
	}
	return backoff.Fibonacci(factor), nil
}

// transforms

func generatedFullTransformation(_ string) (jitter.Transformation, error) {
	return jitter.Full(generator()), nil
}

func generatedEqualTransformation(_ string) (jitter.Transformation, error) {
	return jitter.Equal(generator()), nil
}

func generatedDeviationTransformation(raw string) (jitter.Transformation, error) {
	factor, err := strconv.ParseFloat(raw, 0)
	if err != nil {
		return nil, err
	}
	return jitter.Deviation(generator(), factor), nil
}

func generatedNormalDistributionTransformation(raw string) (jitter.Transformation, error) {
	standardDeviation, err := strconv.ParseFloat(raw, 0)
	if err != nil {
		return nil, err
	}
	return jitter.NormalDistribution(generator(), standardDeviation), nil
}

func init() {
	var (
		fInfinite                                  bool
		fLimit, fDelay, fWait, fBackoff, fTBackoff string
	)
	compliance = map[string]struct {
		cursor  interface{}
		usage   string
		handler func(*flag.Flag) (strategy.Strategy, error)
	}{
		"infinite": {cursor: &fInfinite,
			handler: generatedInfiniteStrategy},
		"limit": {cursor: &fLimit,
			handler: generatedLimitStrategy},
		"delay": {cursor: &fDelay,
			handler: generatedDelayStrategy},
		"wait": {cursor: &fWait,
			handler: generatedWaitStrategy},
		"backoff": {cursor: &fBackoff,
			handler: generatedBackoffStrategy},
		"tbackoff": {cursor: &fTBackoff,
			handler: generatedBackoffWithJitterStrategy},
	}
	algorithms = map[string]func(args string) (backoff.Algorithm, error){
		"inc":    generatedIncrementalAlgorithm,
		"lin":    generatedLinearAlgorithm,
		"epx":    generatedExponentialAlgorithm,
		"binexp": generatedBinaryExponentialAlgorithm,
		"fib":    generatedFibonacciAlgorithm,
	}
	transforms = map[string]func(args string) (jitter.Transformation, error){
		"full":  generatedFullTransformation,
		"equal": generatedEqualTransformation,
		"dev":   generatedDeviationTransformation,
		"ndist": generatedNormalDistributionTransformation,
	}
}
