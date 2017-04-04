package main

import (
	"os"
	"os/exec"

	"github.com/kamilsk/retrier"
)

func main() {
	ctx, args, strategies := parse()
	action := func(attempt uint) error {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return cmd.Run()
	}
	retrier.Retry(ctx, action, strategies...)
}
