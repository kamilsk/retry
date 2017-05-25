package main

import (
	"bytes"
	"log"
	"testing"
)

func TestLogger(t *testing.T) {
	b := bytes.NewBuffer(nil)
	p := log.New(b, "", 0)
	l := logger{stderr: p, stdout: p}
	defaultError := func(t *testing.T, expected, obtained, name string, index int) {
		t.Errorf("unexpected buffer value, expected %q, obtained %q at {%d:%s} test case",
			expected, obtained, index, name)
	}

	for i, tc := range []struct {
		name           string
		action         func()
		expected       string
		debug, noColor bool
		error          func(t *testing.T, expected, obtained, name string, index int)
	}{
		{
			name:     "disabled debug",
			action:   func() { l.Debug("test") },
			expected: "",
			error: func(t *testing.T, expected, obtained, name string, index int) {
				t.Errorf("unexpected buffer value, expected empty string, obtained %q at {%d:%s}",
					obtained, index, name)
			},
		},
		{
			name:     "colored debug",
			action:   func() { l.Debug("logger") },
			expected: "\x1b[33m[DEBUG]\x1b[0m logger\n",
			debug:    true,
			error:    defaultError,
		},
		{
			name:     "not colored debug",
			action:   func() { l.Debug("logger") },
			expected: "[DEBUG] logger\n",
			debug:    true,
			noColor:  true,
			error:    defaultError,
		},
		{
			name:     "disabled debug",
			action:   func() { l.Debugf("%s", "test") },
			expected: "",
			error: func(t *testing.T, expected, obtained, name string, index int) {
				t.Errorf("unexpected buffer value, expected empty string, obtained %q at {%d:%s}",
					obtained, index, name)
			},
		},
		{
			name:     "colored debug",
			action:   func() { l.Debugf("%s", "logger") },
			expected: "\x1b[33m[DEBUG]\x1b[0m logger\n",
			debug:    true,
			error:    defaultError,
		},
		{
			name:     "not colored debug",
			action:   func() { l.Debugf("%s", "logger") },
			expected: "[DEBUG] logger\n",
			debug:    true,
			noColor:  true,
			error:    defaultError,
		},
		{
			name:     "colored error",
			action:   func() { l.Error("logger") },
			expected: "\x1b[31m[ERROR]\x1b[0m logger\n",
			error:    defaultError,
		},
		{
			name:     "not colored error",
			action:   func() { l.Error("logger") },
			expected: "[ERROR] logger\n",
			noColor:  true,
			error:    defaultError,
		},
		{
			name:     "colored error",
			action:   func() { l.Errorf("%s", "logger") },
			expected: "\x1b[31m[ERROR]\x1b[0m logger\n",
			error:    defaultError,
		},
		{
			name:     "not colored error",
			action:   func() { l.Errorf("%s", "logger") },
			expected: "[ERROR] logger\n",
			noColor:  true,
			error:    defaultError,
		},
		{
			name:     "colored info",
			action:   func() { l.Info("logger") },
			expected: "\x1b[32m[INFO]\x1b[0m logger\n",
			error:    defaultError,
		},
		{
			name:     "not colored info",
			action:   func() { l.Info("logger") },
			expected: "[INFO] logger\n",
			noColor:  true,
			error:    defaultError,
		},
		{
			name:     "colored info",
			action:   func() { l.Infof("%s", "logger") },
			expected: "\x1b[32m[INFO]\x1b[0m logger\n",
			error:    defaultError,
		},
		{
			name:     "not colored info",
			action:   func() { l.Infof("%s", "logger") },
			expected: "[INFO] logger\n",
			noColor:  true,
			error:    defaultError,
		},
	} {
		b.Reset()
		l.debug, l.colored = tc.debug, !tc.noColor

		tc.action()

		obtained := b.String()
		if tc.expected != obtained {
			tc.error(t, tc.expected, obtained, tc.name, i)
		}
	}
}
