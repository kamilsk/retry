package main

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
		l.stdout.Print("[DEBUG] " + message)
	}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if l.debug {
		l.stdout.Printf("[DEBUG] "+format, args...)
	}
}

func (l *logger) Error(message string) {
	l.stderr.Print("[ERROR] " + message)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.stderr.Printf("[ERROR] "+format, args...)
}

func (l *logger) Info(message string) {
	l.stdout.Print("[INFO] " + message)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.stdout.Printf("[INFO] "+format, args...)
}
