//+build go1.11

package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCompletion(t *testing.T) {
	before := completionCommand.OutOrStdout()
	defer completionCommand.SetOutput(before)

	buf := bytes.NewBuffer(nil)
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(completionCommand)
	cmd.SetOutput(buf)

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{"Bash", "bash", "# bash completion for test"},
		{"Zsh", "zsh", "#compdef test"},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			buf.Reset()
			assert.NoError(t, completionCommand.Flag("format").Value.Set(tc.format))
			assert.NoError(t, completionCommand.RunE(completionCommand, nil))
			assert.Contains(t, buf.String(), tc.expected)
		})
	}
}
