package main

import "strconv"

const (
	escape = "\x1b"
	red    = iota + 30
	green
	yellow
)

type printer interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

type logger struct {
	stderr printer
	stdout printer
	debug  bool
}

func (l *logger) Debug(message string) {
	if l.debug {
		l.stdout.Print(colorize(yellow, "[DEBUG]") + " " + message)
	}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if l.debug {
		l.stdout.Printf(colorize(yellow, "[DEBUG]")+" "+format, args...)
	}
}

func (l *logger) Error(message string) {
	l.stderr.Print(colorize(red, "[ERROR]") + " " + message)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.stderr.Printf(colorize(red, "[ERROR]")+" "+format, args...)
}

func (l *logger) Info(message string) {
	l.stdout.Print(colorize(green, "[INFO]") + " " + message)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.stdout.Printf(colorize(green, "[INFO]")+" "+format, args...)
}

func colorize(color int, str string) string {
	if NoColor {
		return str
	}
	return escape + "[" + strconv.Itoa(color) + "m" + str + escape + "[0m"
}
