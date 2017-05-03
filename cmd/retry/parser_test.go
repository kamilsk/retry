package main

import (
	"errors"
	"flag"
	"fmt"
	"testing"

	"github.com/kamilsk/retry/strategy"
)

func safe(arguments ...string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case string:
				err = errors.New(r.(string))
			case error:
				err = v
			default:
				err = fmt.Errorf("unexpected panic type %T", v)
			}
		}
	}()

	parse(arguments...)

	return
}

func Test_parse(t *testing.T) {
	for i, tc := range []struct {
		name   string
		before func()
		do     func() (obtained, expected string)
		after  func()
	}{
		{
			name: "unsupported cursor",
			before: func() {
				var unsupported int
				compliance["unsupported"] = struct {
					cursor  interface{}
					usage   string
					handler func(*flag.Flag) (strategy.Strategy, error)
				}{
					cursor: &unsupported,
				}
			},
			do: func() (obtained, expected string) {
				expected = "an unsupported cursor type *int"
				if err := safe(); err != nil {
					obtained = err.Error()
				}
				return
			},
			after: func() {
				delete(compliance, "unsupported")
			},
		},
		{
			name: "invalid arguments",
			do: func() (obtained, expected string) {
				expected = "flag provided but not defined: -test"
				if err := safe("-test=invalid"); err != nil {
					obtained = err.Error()
				}
				return
			},
		},
		{
			name: "invalid timeout",
			do: func() (obtained, expected string) {
				expected = "time: invalid duration timeout"
				if err := safe("-timeout=timeout"); err != nil {
					obtained = err.Error()
				}
				return
			},
		},
		{
			name: "invalid strategy",
			do: func() (obtained, expected string) {
				expected = "time: invalid duration duration"
				if err := safe("-delay=duration"); err != nil {
					obtained = err.Error()
				}
				return
			},
		},
		{
			name: "nothing to do",
			do: func() (obtained, expected string) {
				expected = "please provide a command to retry"
				if err := safe("-delay=1s"); err != nil {
					obtained = err.Error()
				}
				return
			},
		},
		{
			name: "success",
			do: func() (obtained, expected string) {
				expected = ""
				if err := safe("-delay=1s", "--", "whoami"); err != nil {
					obtained = err.Error()
				}
				return
			},
		},
	} {
		if tc.before != nil {
			tc.before()
		}

		if obtained, expected := tc.do(); obtained != expected {
			t.Errorf("expected panic with message %q, obtained %q at {%s:%d}", expected, obtained, tc.name, i)
		}

		if tc.after != nil {
			tc.after()
		}
	}
}

func Test_handle(t *testing.T) {
	for i, tc := range []struct {
		name     string
		flags    []*flag.Flag
		error    string
		expected int
	}{
		{
			name:     "empty list",
			flags:    nil,
			expected: 0,
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

func Test_parseAlgorithm(t *testing.T) {
	for i, tc := range []struct {
		name  string
		args  string
		error string
	}{
		{
			name:  "not matched",
			args:  "~",
			error: "invalid argument ~",
		},
		{
			name:  "unknown algorithm",
			args:  "x",
			error: "unknown algorithm x",
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

func Test_parseTransform(t *testing.T) {
	for i, tc := range []struct {
		name  string
		args  string
		error string
	}{
		{
			name:  "not matched",
			args:  "~",
			error: "invalid argument ~",
		},
		{
			name:  "unknown transformation",
			args:  "x",
			error: "unknown transformation x",
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
