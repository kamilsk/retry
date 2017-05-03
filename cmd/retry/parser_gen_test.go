package main

// TODO:GEN generate it

import (
	"bytes"
	"flag"
	"fmt"
	"testing"
)

type value string

func (v value) String() string {
	return string(v)
}

func (v value) Set(n string) error {
	p := &v
	*p = value(n)
	return nil
}

func Test_handle_generated(t *testing.T) {
	for i, tc := range []struct {
		name     string
		flags    []*flag.Flag
		error    string
		expected int
	}{
		{
			name: "infinite",
			flags: []*flag.Flag{
				{
					Name: "infinite",
				},
			},
			expected: 1,
		},
		{
			name: "limit",
			flags: []*flag.Flag{
				{
					Name:  "limit",
					Value: value("1"),
				},
			},
			expected: 1,
		},
		{
			name: "limit: invalid attemptLimit",
			flags: []*flag.Flag{
				{
					Name:  "limit",
					Value: value("attemptLimit"),
				},
			},
			error: `strconv.ParseUint: parsing "attemptLimit": invalid syntax`,
		},
		{
			name: "delay",
			flags: []*flag.Flag{
				{
					Name:  "delay",
					Value: value("1s"),
				},
			},
			expected: 1,
		},
		{
			name: "delay: invalid duration",
			flags: []*flag.Flag{
				{
					Name:  "delay",
					Value: value("duration"),
				},
			},
			error: "time: invalid duration duration",
		},
		{
			name: "wait",
			flags: []*flag.Flag{
				{
					Name:  "wait",
					Value: value("1s,1s,1s,1s,1s"),
				},
			},
			expected: 1,
		},
		{
			name: "wait: invalid duration",
			flags: []*flag.Flag{
				{
					Name:  "wait",
					Value: value("1s,1s,duration,1s,1s"),
				},
			},
			error: "time: invalid duration duration",
		},
		{
			name: "backoff",
			flags: []*flag.Flag{
				{
					Name:  "backoff",
					Value: value("inc[1s,1s]"),
				},
			},
			expected: 1,
		},
		{
			name: "backoff: unknown algorithm",
			flags: []*flag.Flag{
				{
					Name:  "backoff",
					Value: value("x"),
				},
			},
			error: "unknown algorithm x",
		},
		{
			name: "backoff with jitter",
			flags: []*flag.Flag{
				{
					Name:  "tbackoff",
					Value: value("inc[1s,1s] full"),
				},
			},
			expected: 1,
		},
		{
			name: "backoff with jitter: invalid argument count",
			flags: []*flag.Flag{
				{
					Name:  "tbackoff",
					Value: value("inc[1s,1s]"),
				},
			},
			error: "invalid argument count",
		},
		{
			name: "backoff with jitter: unknown algorithm",
			flags: []*flag.Flag{
				{
					Name:  "tbackoff",
					Value: value("x full"),
				},
			},
			error: "unknown algorithm x",
		},
		{
			name: "backoff with jitter: unknown transformation",
			flags: []*flag.Flag{
				{
					Name:  "tbackoff",
					Value: value("inc[1s,1s] x"),
				},
			},
			error: "unknown transformation x",
		},
	} {
		strategies, err := handle(tc.flags)

		switch {
		case tc.error == "" && err != nil:
			t.Errorf("unexpected error %q at {%s:%d} test case", err, tc.name, i)
		case tc.error != "" && err == nil:
			t.Errorf("expected error %q, obtained nil at {%s:%d} test case", tc.error, tc.name, i)
		case tc.error != "" && err != nil:
			if tc.error != err.Error() {
				t.Errorf("expected error %q, obtained %q at {%s:%d} test case", tc.error, err, tc.name, i)
			}
		case len(strategies) != tc.expected:
			t.Errorf("expected %d strategies, obtained %d", tc.expected, len(strategies))
		}
	}
}

func Test_parseAlgorithm_generated(t *testing.T) {
	for i, tc := range []struct {
		name  string
		args  string
		error string
	}{
		{
			name: "incremental",
			args: "inc[1s,1s]",
		},
		{
			name:  "incremental: invalid argument count",
			args:  "inc[1s]",
			error: "invalid argument count",
		},
		{
			name:  "incremental: invalid initial",
			args:  "inc[initial,1s]",
			error: "time: invalid duration initial",
		},
		{
			name:  "incremental: invalid increment",
			args:  "inc[1s,increment]",
			error: "time: invalid duration increment",
		},
		{
			name: "linear",
			args: "lin[1s]",
		},
		{
			name:  "linear: invalid factor",
			args:  "lin[factor]",
			error: "time: invalid duration factor",
		},
		{
			name: "exponential",
			args: "exp[1s,1.0]",
		},
		{
			name:  "exponential: invalid argument count",
			args:  "exp[1s]",
			error: "invalid argument count",
		},
		{
			name:  "exponential: invalid factor",
			args:  "exp[factor,1.0]",
			error: "time: invalid duration factor",
		},
		{
			name:  "exponential: invalid base",
			args:  "exp[1s,base]",
			error: `strconv.ParseFloat: parsing "base": invalid syntax`,
		},
		{
			name: "binary exponential",
			args: "binexp[1s]",
		},
		{
			name:  "binary exponential: invalid factor",
			args:  "binexp[factor]",
			error: "time: invalid duration factor",
		},
		{
			name: "fibonacci",
			args: "fib[1s]",
		},
		{
			name:  "fibonacci: invalid factor",
			args:  "fib[factor]",
			error: "time: invalid duration factor",
		},
	} {
		alg, err := parseAlgorithm(tc.args)

		switch {
		case tc.error == "" && err != nil:
			t.Errorf("unexpected error %q at {%s:%d} test case", err, tc.name, i)
		case tc.error != "" && err == nil:
			t.Errorf("expected error %q, obtained nil at {%s:%d} test case", tc.error, tc.name, i)
		case tc.error != "" && err != nil:
			if tc.error != err.Error() {
				t.Errorf("expected error %q, obtained %q at {%s:%d} test case", tc.error, err, tc.name, i)
			}
		case alg == nil:
			t.Error("expected correct algorithm, obtained nil")
		}
	}
}

func Test_parseTransform_generated(t *testing.T) {
	for i, tc := range []struct {
		name  string
		args  string
		error string
	}{
		{
			name: "full",
			args: "full",
		},
		{
			name: "equal",
			args: "equal",
		},
		{
			name: "deviation",
			args: "dev[1.0]",
		},
		{
			name:  "deviation: invalid factor",
			args:  "dev[factor]",
			error: `strconv.ParseFloat: parsing "factor": invalid syntax`,
		},
		{
			name: "normal distribution",
			args: "ndist[1.0]",
		},
		{
			name:  "normal distribution: invalid standardDeviation",
			args:  "ndist[standardDeviation]",
			error: `strconv.ParseFloat: parsing "standardDeviation": invalid syntax`,
		},
	} {
		tr, err := parseTransform(tc.args)

		switch {
		case tc.error == "" && err != nil:
			t.Errorf("unexpected error %q at {%s:%d} test case", err, tc.name, i)
		case tc.error != "" && err == nil:
			t.Errorf("expected error %q, obtained nil at {%s:%d} test case", tc.error, tc.name, i)
		case tc.error != "" && err != nil:
			if tc.error != err.Error() {
				t.Errorf("expected error %q, obtained %q at {%s:%d} test case", tc.error, err, tc.name, i)
			}
		case tr == nil:
			t.Error("expected correct transformation, obtained nil")
		}
	}
}

func Test_usage(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	// don't forget update README.md
	expected := fmt.Sprintf(`
usage: test [-timeout timeout] [strategy flags] -- command

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
    inc[Xs,Ys]
        Incremental creates a Algorithm that increments the initial duration
        by the given increment for each attempt.
    lin[Xs]
        Linear creates a Algorithm that linearly multiplies the factor
        duration by the attempt number for each attempt.
    exp[Xs,Y]
        Exponential creates a Algorithm that multiplies the factor duration by
        an exponentially increasing factor for each attempt, where the factor is
        calculated as the given base raised to the attempt number.
    binexp[Xs]
        BinaryExponential creates a Algorithm that multiplies the factor
        duration by an exponentially increasing factor for each attempt, where the
        factor is calculated as "2" raised to the attempt number (2^attempt).
    fib[Xs]
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
    dev[X]
        Deviation creates a Transformation that transforms a duration into a result
        duration that deviates from the input randomly by a given factor.

        The given generator is what is used to determine the random transformation.
        If a nil generator is passed, a default one will be provided.

        Inspired by https://developers.google.com/api-client-library/java/google-http-java-client/backoff
    ndist[X]
        NormalDistribution creates a Transformation that transforms a duration into a
        result duration based on a normal distribution of the input and the given
        standard deviation.

        The given generator is what is used to determine the random transformation.
        If a nil generator is passed, a default one will be provided.

Full example:
    retry -limit=3 -backoff=lin[10ms] -- curl http://unknown.host
    retry -timeout=500ms --infinite -- curl http://unknown.host

Current version is %s.
`, Version)

	usage(buf, "test")

	if buf.String() != expected {
		t.Error("unexpected usage message")
	}
}
