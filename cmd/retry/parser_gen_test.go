package main

import "testing"

// TODO:GEN generate it

func Test_handle_generated(t *testing.T) {
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
			name:  "incremental: argument count",
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
			name:  "exponential: argument count",
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
			args:  "exp[1s,1s]",
			error: `strconv.ParseFloat: parsing "1s": invalid syntax`,
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
}
