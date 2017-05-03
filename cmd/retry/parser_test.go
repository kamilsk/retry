package main

import (
	"flag"
	"testing"
)

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
