package main

import (
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/kamilsk/retry"
)

const (
	Success = 0
	Failed  = 1
)

func main() { application{Args: os.Args, Stderr: os.Stderr, Stdout: os.Stdout, Shutdown: os.Exit}.Run() }

type application struct {
	Args           []string
	Stderr, Stdout io.Writer
	Shutdown       func(code int)
}

// Run executes the application logic.
func (app application) Run() {
	var (
		start   time.Time
		started bool
	)

	result, err := parse(app.Stderr, app.Args[0], app.Args[1:]...)
	format := "an error occurred: %v\n"
	if result.Debug {
		format = "an error occurred: %+v\n"
	}
	if err != nil {
		color.New(color.FgRed).Fprintf(app.Stderr, format, err)
		app.Shutdown(Failed)
		return
	}
	defer func() {
		if result.Notify {
			// TODO try to find or implement by myself
			// - https://github.com/variadico/noti
			// - https://github.com/jolicode/JoliNotif
			color.New(color.FgYellow).Fprintln(app.Stderr, "notify component is not ready yet")
		}
	}()

	action := func(attempt uint) error {
		if !started {
			start = time.Now()
			started = true
		} else {
			color.New(color.FgYellow).Fprintf(app.Stderr, "#%d attempt at %s... \n", attempt+1, time.Now().Sub(start))
		}
		cmd := exec.Command(result.Args[0], result.Args[1:]...)
		cmd.Stdout, cmd.Stderr = app.Stdout, app.Stderr
		return cmd.Run()
	}
	deadline := retry.Multiplex(
		retry.WithTimeout(result.Timeout),
		retry.WithSignal(os.Interrupt),
	)
	if err := retry.Retry(deadline, action, result.Strategies...); err != nil {
		color.New(color.FgRed).Fprintf(app.Stderr, format, err)
		app.Shutdown(Failed)
		return
	}
	app.Shutdown(Success)
	return
}
