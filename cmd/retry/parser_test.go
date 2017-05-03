package main

import (
	"testing"
)

func Test_handle(t *testing.T) {
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
}
