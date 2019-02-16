package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	type memo struct {
		commit  string
		date    string
		version string
	}

	buf := bytes.NewBuffer(nil)
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(versionCommand)
	cmd.SetOutput(buf)

	tests := []struct {
		name string
		memo
		expected string
	}{
		{"Version 4.0", memo{commit: "...", date: "Sometime", version: "4.0.0"},
			"Version 4.0.0 (commit: ..., build date: Sometime"},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			before := memo{commit: commit, date: date, version: version}
			defer func() { commit, date, version = before.commit, before.date, before.version }()
			commit, date, version = tc.commit, tc.date, tc.version

			buf.Reset()
			versionCommand.Run(versionCommand, nil)
			assert.Contains(t, buf.String(), tc.expected)
		})
	}
}
