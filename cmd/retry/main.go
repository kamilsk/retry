package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/kamilsk/retry"
)

var (
	// Debug prints verbose information to stdout.
	// Can be changed by `-ldflags "-X" 'main.Debug=..."'`
	// or `-v` parameter.
	Debug = false
	// NoColor deprecates colorize logger' output.
	// Can be changed by `-ldflags "-X 'main.NoColor=...'"`.
	NoColor = false
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error occurred %q \n", r)
			os.Exit(1)
		}
	}()

	var stderrFlag, stdoutFlag int
	if Debug {
		stderrFlag = log.Lshortfile
	}

	l := &logger{
		stderr:  log.New(os.Stderr, "", stderrFlag),
		stdout:  log.New(os.Stdout, "", stdoutFlag),
		debug:   Debug,
		colored: !NoColor,
	}

	var (
		start   time.Time
		started bool
	)

	result, err := parse(os.Args[0], os.Args[1:]...)
	if err != nil {
		l.Errorf("error occurred: %q", err)
		os.Exit(1)
	}
	defer func() {
		if result.Notify {
			// TODO try to find or implement by myself
			// - https://github.com/variadico/noti
			// - https://github.com/jolicode/JoliNotif
			color.New(color.FgYellow).Fprintln(os.Stderr, "notify component is not ready yet")
		}
	}()

	action := func(attempt uint) error {
		if !started {
			start = time.Now()
			started = true
		} else {
			l.Infof("#%d attempt at %s... \n", attempt+1, time.Now().Sub(start))
		}
		cmd := exec.Command(result.Args[0], result.Args[1:]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return cmd.Run()
	}
	deadline := retry.Multiplex(
		retry.WithTimeout(result.Timeout),
		retry.WithSignal(os.Interrupt),
	)
	if err := retry.Retry(deadline, action, result.Strategies...); err != nil {
		l.Errorf("error occurred: %q", err)
		os.Exit(1)
	}
}
