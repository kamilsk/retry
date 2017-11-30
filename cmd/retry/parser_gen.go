package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/kamilsk/retry/backoff"
	"github.com/kamilsk/retry/jitter"
	"github.com/kamilsk/retry/strategy"
)

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
		"exp":    generatedExponentialAlgorithm,
		"binexp": generatedBinaryExponentialAlgorithm,
		"fib":    generatedFibonacciAlgorithm,
	}
	transforms = map[string]func(args string) (jitter.Transformation, error){
		"full":  generatedFullTransformation,
		"equal": generatedEqualTransformation,
		"dev":   generatedDeviationTransformation,
		"ndist": generatedNormalDistributionTransformation,
	}
	usage = func(output io.Writer, md Metadata) func() {
		return func() {
			fmt.Fprintf(output, `
Usage: %s [-timeout Timeout] [--debug] [--notify] [strategy flags] -- command

The strategy flags
    --infinite
        Infinite creates a Strategy that will never stop repeating.

    -limit=X
        Limit creates a Strategy that limits the number of attempts that Retry will
        make.

    -delay=Xs
        Delay creates a Strategy that waits the given duration before the first
        attempt is made.

    -wait=Xs,...
        Wait creates a Strategy that waits the given durations for each attempt after
        the first. If the number of attempts is greater than the number of durations
        provided, then the strategy uses the last duration provided.

    -backoff=:algorithm
        Backoff creates a Strategy that waits before each attempt, with a duration as
        defined by the given backoff.Algorithm.

    -tbackoff=":algorithm :transformation"
        BackoffWithJitter creates a Strategy that waits before each attempt, with a
        duration as defined by the given backoff.Algorithm and jitter.Transformation.

:algorithm
    inc:Xs,Ys
        Incremental creates a Algorithm that increments the initial duration
        by the given increment for each attempt.

    lin:Xs
        Linear creates a Algorithm that linearly multiplies the factor
        duration by the attempt number for each attempt.

    exp:Xs,Y
        Exponential creates a Algorithm that multiplies the factor duration by
        an exponentially increasing factor for each attempt, where the factor is
        calculated as the given base raised to the attempt number.

    binexp:Xs
        BinaryExponential creates a Algorithm that multiplies the factor
        duration by an exponentially increasing factor for each attempt, where the
        factor is calculated as "2" raised to the attempt number (2^attempt).

    fib:Xs
        Fibonacci creates a Algorithm that multiplies the factor duration by
        an increasing factor for each attempt, where the factor is the Nth number in
        the Fibonacci sequence.

:transformation
    full
        Full creates a Transformation that transforms a duration into a result
        duration in [0, n) randomly, where n is the given duration.

        The given generator is what is used to determine the random transformation.
        If a nil generator is passed, a default one will be provided.

        Inspired by https://www.awsarchitectureblog.com/2015/03/backoff.html

    equal
        Equal creates a Transformation that transforms a duration into a result
        duration in [n/2, n) randomly, where n is the given duration.

        The given generator is what is used to determine the random transformation.
        If a nil generator is passed, a default one will be provided.

        Inspired by https://www.awsarchitectureblog.com/2015/03/backoff.html

    dev:X
        Deviation creates a Transformation that transforms a duration into a result
        duration that deviates from the input randomly by a given factor.

        The given generator is what is used to determine the random transformation.
        If a nil generator is passed, a default one will be provided.

        Inspired by https://developers.google.com/api-client-library/java/google-http-java-client/backoff

    ndist:X
        NormalDistribution creates a Transformation that transforms a duration into a
        result duration based on a normal distribution of the input and the given
        standard deviation.

        The given generator is what is used to determine the random transformation.
        If a nil generator is passed, a default one will be provided.

Examples:
    %[1]s -limit=3 -backoff=lin:10ms -- curl http://example.com
    %[1]s -tbackoff="lin:10s full" -- curl https://example.com
    %[1]s -timeout=500ms --notify --infinite -- git pull

Version %s (commit: %s, build date: %s, go version: %s, compiler: %s, platform: %s)
`, md.BinName, md.Version, md.Commit, md.BuildDate, md.GoVersion, md.Compiler, md.Platform)
		}
	}
}

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
	args := strings.Split(f.Value.String(), " ")
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
	increment, err := time.ParseDuration(args[1])
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
	return jitter.Full(rand.New(rand.NewSource(time.Now().UnixNano()))), nil
}

func generatedEqualTransformation(_ string) (jitter.Transformation, error) {
	return jitter.Equal(rand.New(rand.NewSource(time.Now().UnixNano()))), nil
}

func generatedDeviationTransformation(raw string) (jitter.Transformation, error) {
	factor, err := strconv.ParseFloat(raw, 0)
	if err != nil {
		return nil, err
	}
	return jitter.Deviation(rand.New(rand.NewSource(time.Now().UnixNano())), factor), nil
}

func generatedNormalDistributionTransformation(raw string) (jitter.Transformation, error) {
	standardDeviation, err := strconv.ParseFloat(raw, 0)
	if err != nil {
		return nil, err
	}
	return jitter.NormalDistribution(rand.New(rand.NewSource(time.Now().UnixNano())), standardDeviation), nil
}
