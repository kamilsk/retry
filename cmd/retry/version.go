// +build go1.10

package main

import (
	"runtime"

	"github.com/spf13/cobra"
)

var (
	commit  = "none"
	date    = "unknown"
	version = "dev"
)

// Version shows application version.
var Version = &cobra.Command{
	Use:   "version",
	Short: "Show application version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf(
			"Version %s (commit: %s, build date: %s, go version: %s, compiler: %s, platform: %s/%s)\n",
			version, commit, date, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
	},
	Version: version,
}
