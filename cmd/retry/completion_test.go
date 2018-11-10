// +build go1.10

package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompletion(t *testing.T) {
	before := Completion.OutOrStdout()
	defer Completion.SetOutput(before)

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{"Bash", "bash", "# bash completion for completion"},
		{"Zsh", "zsh", "#compdef completion"},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			Completion.SetOutput(buf)
			Completion.Flag("format").Value.Set(tc.format)
			assert.NoError(t, Completion.RunE(Completion, nil))
			assert.Contains(t, buf.String(), tc.expected)
		})
	}
}
