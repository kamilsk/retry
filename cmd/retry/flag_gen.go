package main

import (
	"errors"
	"flag"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kamilsk/retrier/backoff"
	"github.com/kamilsk/retrier/jitter"
	"github.com/kamilsk/retrier/strategy"
)

type Compliance map[string]struct {
	cursor  interface{}
	usage   string
	handler func(*flag.Flag) (strategy.Strategy, error)
}

var (
	compliance Compliance
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

func InfiniteStrategy_gen(_ *flag.Flag) (strategy.Strategy, error) {
	return strategy.Infinite(), nil
}

func LimitStrategy_gen(f *flag.Flag) (strategy.Strategy, error) {
	attemptLimit, err := strconv.ParseUint(f.Value.String(), 10, 0)
	if err != nil {
		return nil, err
	}
	return strategy.Limit(uint(attemptLimit)), nil
}

func DelayStrategy_gen(f *flag.Flag) (strategy.Strategy, error) {
	duration, err := time.ParseDuration(f.Value.String())
	if err != nil {
		return nil, err
	}
	return strategy.Delay(duration), nil
}

func WaitStrategy_gen(f *flag.Flag) (strategy.Strategy, error) {
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

func BackoffStrategy_gen(f *flag.Flag) (strategy.Strategy, error) {
	algorithm, err := parseAlgorithm(f.Value.String())
	if err != nil {
		return nil, err
	}
	return strategy.Backoff(algorithm), nil
}

func BackoffWithJitterStrategy_gen(f *flag.Flag) (strategy.Strategy, error) {
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

func IncrementalAlgorithm_gen(raw string) (backoff.Algorithm, error) {
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

func LinearAlgorithm_gen(raw string) (backoff.Algorithm, error) {
	factor, err := time.ParseDuration(raw)
	if err != nil {
		return nil, err
	}
	return backoff.Linear(factor), nil
}

func ExponentialAlgorithm_gen(raw string) (backoff.Algorithm, error) {
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

func BinaryExponentialAlgorithm_gen(raw string) (backoff.Algorithm, error) {
	factor, err := time.ParseDuration(raw)
	if err != nil {
		return nil, err
	}
	return backoff.BinaryExponential(factor), nil
}

func FibonacciAlgorithm_gen(raw string) (backoff.Algorithm, error) {
	factor, err := time.ParseDuration(raw)
	if err != nil {
		return nil, err
	}
	return backoff.Fibonacci(factor), nil
}

// transforms

func FullTransformation_gen(_ string) (jitter.Transformation, error) {
	return jitter.Full(generator()), nil
}

func EqualTransformation_gen(_ string) (jitter.Transformation, error) {
	return jitter.Equal(generator()), nil
}

func DeviationTransformation_gen(raw string) (jitter.Transformation, error) {
	factor, err := strconv.ParseFloat(raw, 0)
	if err != nil {
		return nil, err
	}
	return jitter.Deviation(generator(), factor), nil
}

func NormalDistributionTransformation_gen(raw string) (jitter.Transformation, error) {
	standardDeviation, err := strconv.ParseFloat(raw, 0)
	if err != nil {
		return nil, err
	}
	return jitter.NormalDistribution(generator(), standardDeviation), nil
}

func init() {
	var (
		f_infinite                                      bool
		f_limit, f_delay, f_wait, f_backoff, f_tbackoff string
	)
	compliance = Compliance{
		"infinite": {cursor: &f_infinite,
			handler: InfiniteStrategy_gen},
		"limit": {cursor: &f_limit,
			handler: LimitStrategy_gen},
		"delay": {cursor: &f_delay,
			handler: DelayStrategy_gen},
		"wait": {cursor: &f_wait,
			handler: WaitStrategy_gen},
		"backoff": {cursor: &f_backoff,
			handler: BackoffStrategy_gen},
		"tbackoff": {cursor: &f_tbackoff,
			handler: BackoffWithJitterStrategy_gen},
	}
	algorithms = map[string]func(args string) (backoff.Algorithm, error){
		"inc":    IncrementalAlgorithm_gen,
		"lin":    LinearAlgorithm_gen,
		"epx":    ExponentialAlgorithm_gen,
		"binexp": BinaryExponentialAlgorithm_gen,
		"fib":    FibonacciAlgorithm_gen,
	}
	transforms = map[string]func(args string) (jitter.Transformation, error){
		"full":  FullTransformation_gen,
		"equal": EqualTransformation_gen,
		"dev":   DeviationTransformation_gen,
		"ndist": NormalDistributionTransformation_gen,
	}
}
