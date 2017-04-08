package main

import (
	"flag"
	"math/rand"
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
)

func gen() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// TODO: generate it

// strategies

func InfiniteStrategy_gen(_ *flag.Flag) (strategy.Strategy, error) {
	return strategy.Infinite(), nil
}

func LimitStrategy_gen(f *flag.Flag) (strategy.Strategy, error) {
	return strategy.Limit(0), nil
}

func DelayStrategy_gen(f *flag.Flag) (strategy.Strategy, error) {
	return strategy.Delay(0), nil
}

func WaitStrategy_gen(f *flag.Flag) (strategy.Strategy, error) {
	return strategy.Wait(0, 0, 0), nil
}

func BackoffStrategy_gen(f *flag.Flag) (strategy.Strategy, error) {
	return strategy.Backoff(nil), nil
}

func BackoffWithJitterStrategy_gen(f *flag.Flag) (strategy.Strategy, error) {
	return strategy.BackoffWithJitter(nil, nil), nil
}

// algorithms

func IncrementalAlgorithm_gen(args string) (backoff.Algorithm, error) {
	return backoff.Incremental(0, 0), nil
}

func LinearAlgorithm_gen(args string) (backoff.Algorithm, error) {
	return backoff.Linear(0), nil
}

func ExponentialAlgorithm_gen(args string) (backoff.Algorithm, error) {
	return backoff.Exponential(0, 0), nil
}

func BinaryExponentialAlgorithm_gen(args string) (backoff.Algorithm, error) {
	return backoff.BinaryExponential(0), nil
}

func FibonacciAlgorithm_gen(args string) (backoff.Algorithm, error) {
	return backoff.Fibonacci(0), nil
}

// transforms

func FullTransformation_gen(_ string) (jitter.Transformation, error) {
	return jitter.Full(gen()), nil
}

func EqualTransformation_gen(_ string) (jitter.Transformation, error) {
	return jitter.Equal(gen()), nil
}

func DeviationTransformation_gen(args string) (jitter.Transformation, error) {
	return jitter.Deviation(gen(), 0), nil
}

func NormalDistributionTransformation_gen(args string) (jitter.Transformation, error) {
	return jitter.NormalDistribution(gen(), 0), nil
}

func init() {
	var (
		f_infinite                              bool
		f_limit, f_delay, f_backoff, f_tbackoff string
	)
	compliance = Compliance{
		"infinite": {cursor: &f_infinite,
			handler: InfiniteStrategy_gen},
		"limit": {cursor: &f_limit,
			handler: LimitStrategy_gen},
		"delay": {cursor: &f_delay,
			handler: DelayStrategy_gen},
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
