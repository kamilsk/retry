package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kamilsk/retry"
)

// Timeout is a timeout of retried operation.
// Can be changed by `-ldflags "-X 'main.Timeout=...'"` or `-timeout ...` parameter.
var Timeout = "1m"

func main() {
	ctx, args, strategies := parse()
	action := func(attempt uint) error {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return cmd.Run()
	}
	if err := retry.Retry(ctx, action, strategies...); err != nil {
		fmt.Fprintf(os.Stderr, "error occurred %q \n", err)
	}
}
