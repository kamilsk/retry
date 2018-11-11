// +build go1.10

package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCompletion(t *testing.T) {
	before := Completion.OutOrStdout()
	defer Completion.SetOutput(before)

	buf := bytes.NewBuffer(nil)
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(Completion)
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
			Completion.Flag("format").Value.Set(tc.format)
			assert.NoError(t, Completion.RunE(Completion, nil))
			assert.Contains(t, buf.String(), tc.expected)
		})
	}
}
