package flag

// Copyright 2009 The Go Authors. All rights reserved.

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

// -- bool Value
type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

func (b *boolValue) Get() interface{} { return bool(*b) }

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }

// optional interface to indicate boolean flags that can be
// supplied without "=value" text
type boolFlag interface {
	flag.Value
	IsBoolFlag() bool
}

// -- string Value
type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() interface{} { return string(*s) }

func (s *stringValue) String() string { return string(*s) }

// A FlagSet represents a set of defined flags. The zero value of a FlagSet
// has no name and has ContinueOnError error handling.
type FlagSet struct {
	// Usage is the function called when an error occurs while parsing flags.
	// The field is a function (not a method) that may be changed to point to
	// a custom error handler.
	Usage func()

	name          string
	parsed        bool
	actual        map[string]*flag.Flag
	formal        map[string]*flag.Flag
	args          []string // arguments after flags
	errorHandling flag.ErrorHandling
	output        io.Writer // nil means stderr; use out() accessor

	sequence []*flag.Flag
}

// Args returns the non-flag arguments.
func (fs *FlagSet) Args() []string { return fs.args }

func (fs *FlagSet) out() io.Writer {
	if fs.output == nil {
		return os.Stderr
	}
	return fs.output
}

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (fs *FlagSet) SetOutput(output io.Writer) {
	fs.output = output
}

// BoolVar defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func (fs *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	fs.Var(newBoolValue(value, p), name, usage)
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func (fs *FlagSet) StringVar(p *string, name string, value string, usage string) {
	fs.Var(newStringValue(value, p), name, usage)
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (fs *FlagSet) Var(value flag.Value, name string, usage string) {
	// Remember the default value as a string; it won't change.
	f := &flag.Flag{Name: name, Usage: usage, Value: value, DefValue: value.String()}
	_, alreadythere := fs.formal[name]
	if alreadythere {
		var msg string
		if fs.name == "" {
			msg = fmt.Sprintf("flag redefined: %s", name)
		} else {
			msg = fmt.Sprintf("%s flag redefined: %s", fs.name, name)
		}
		fmt.Fprintln(fs.out(), msg)
		panic(msg) // Happens only if flags are declared with identical names
	}
	if fs.formal == nil {
		fs.formal = make(map[string]*flag.Flag)
	}
	fs.formal[name] = f
}

// failf prints to standard error a formatted error and
// returns the error.
func (fs *FlagSet) failf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(fs.out(), err)
	fs.usage()
	return err
}

// usage calls the Usage method for the flag set if one is specified.
func (fs *FlagSet) usage() {
	if fs.Usage != nil {
		fs.Usage()
	}
}

// parseOne parses one flag. It reports whether a flag was seen.
func (fs *FlagSet) parseOne() (bool, error) {
	if len(fs.args) == 0 {
		return false, nil
	}
	s := fs.args[0]
	if len(s) == 0 || s[0] != '-' || len(s) == 1 {
		return false, nil
	}
	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags
			fs.args = fs.args[1:]
			return false, nil
		}
	}
	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, fs.failf("bad flag syntax: %s", s)
	}

	// it's a flag. does it have an argument?
	fs.args = fs.args[1:]
	hasValue := false
	value := ""
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
			value = name[i+1:]
			hasValue = true
			name = name[0:i]
			break
		}
	}
	m := fs.formal
	f, alreadythere := m[name] // BUG
	if !alreadythere {
		if name == "help" || name == "h" { // special case for nice help message.
			fs.usage()
			return false, flag.ErrHelp
		}
		return false, fs.failf("flag provided but not defined: -%s", name)
	}

	if fv, ok := f.Value.(boolFlag); ok && fv.IsBoolFlag() { // special case: doesn't need an arg
		if hasValue {
			if err := fv.Set(value); err != nil {
				return false, fs.failf("invalid boolean value %q for -%s: %v", value, name, err)
			}
		} else {
			if err := fv.Set("true"); err != nil {
				return false, fs.failf("invalid boolean flag %s: %v", name, err)
			}
		}
	} else {
		// It must have a value, which might be the next argument.
		if !hasValue && len(fs.args) > 0 {
			// value is the next arg
			hasValue = true
			value, fs.args = fs.args[0], fs.args[1:]
		}
		if !hasValue {
			return false, fs.failf("flag needs an argument: -%s", name)
		}
		if err := f.Value.Set(value); err != nil {
			return false, fs.failf("invalid value %q for flag -%s: %v", value, name, err)
		}
	}
	if fs.actual == nil {
		fs.actual = make(map[string]*flag.Flag)
	}
	fs.actual[name] = f
	fs.sequence = append(fs.sequence, f)
	return true, nil
}

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help or -h were set but not defined.
func (fs *FlagSet) Parse(arguments []string) error {
	fs.parsed = true
	fs.args = arguments
	for {
		seen, err := fs.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}
		switch fs.errorHandling {
		case flag.ContinueOnError:
			return err
		case flag.ExitOnError:
			os.Exit(2)
		case flag.PanicOnError:
			panic(err)
		}
	}
	return nil
}

// NewFlagSet returns a new, empty flag set with the specified name.
func NewFlagSet(name string) *FlagSet {
	return &FlagSet{name: name, errorHandling: flag.PanicOnError}
}

// Flags returns the flag sequence.
func (fs *FlagSet) Flags() []*flag.Flag {
	return fs.sequence
}
