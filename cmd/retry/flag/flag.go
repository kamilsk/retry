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

// A Set represents a set of defined flags. The zero value of a Set
// has no name and has ContinueOnError error handling.
type Set struct {
	// Usage is the function called when an error occurs while parsing flags.
	// The field is a function (not a method) that may be changed to point to
	// a custom error handler.
	Usage func(output io.Writer, args ...string)

	name          string
	formal        map[string]*flag.Flag
	args          []string // arguments after flags
	errorHandling flag.ErrorHandling

	actual []*flag.Flag
}

// Args returns the non-flag arguments.
func (fs *Set) Args() []string { return fs.args }

// BoolVar defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func (fs *Set) BoolVar(p *bool, name string, value bool, usage string) {
	fs.Var(newBoolValue(value, p), name, usage)
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func (fs *Set) StringVar(p *string, name string, value string, usage string) {
	fs.Var(newStringValue(value, p), name, usage)
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (fs *Set) Var(value flag.Value, name string, usage string) {
	// Remember the default value as a string; it won't change.
	f := &flag.Flag{Name: name, Usage: usage, Value: value, DefValue: value.String()}
	if fs.formal == nil {
		fs.formal = make(map[string]*flag.Flag)
	}
	fs.formal[name] = f
}

// failf prints to standard error a formatted error and
// returns the error.
func (fs *Set) failf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

// usage calls the Usage method for the flag set if one is specified.
func (fs *Set) usage() {
	if fs.Usage != nil {
		fs.Usage(os.Stderr, os.Args...)
	}
}

// parseOne parses one flag. It reports whether a flag was seen.
func (fs *Set) parseOne() (bool, error) {
	var (
		result bool
		err    error
	)

	name, err := fs.validate()
	if name == "" {
		return false, err
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

	f, alreadythere := fs.formal[name]
	if !alreadythere {
		if name == "help" || name == "h" { // special case for nice help message.
			fs.usage()
			return false, flag.ErrHelp
		}
		return false, fs.failf("flag provided but not defined: -%s", name)
	}

	if fv, ok := f.Value.(boolFlag); ok && fv.IsBoolFlag() { // special case: doesn't need an arg
		result, err = fs.parseBool(f, name, value, hasValue)
	} else {
		result, err = fs.parseString(f, name, value, hasValue)
	}

	if result {
		fs.actual = append(fs.actual, f)
	}

	return result, err
}

func (fs *Set) parseBool(f *flag.Flag, name, value string, hasValue bool) (bool, error) {
	if hasValue {
		if err := f.Value.Set(value); err != nil {
			return false, fs.failf("invalid boolean value %q for -%s: %v", value, name, err)
		}
	} else {
		if err := f.Value.Set("true"); err != nil {
			return false, fs.failf("invalid boolean flag %s: %v", name, err)
		}
	}
	return true, nil
}

func (fs *Set) parseString(f *flag.Flag, name, value string, hasValue bool) (bool, error) {
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
	return true, nil
}

func (fs *Set) validate() (string, error) {
	if len(fs.args) == 0 {
		return "", nil
	}
	s := fs.args[0]
	if len(s) == 0 || s[0] != '-' || len(s) == 1 {
		return "", nil
	}
	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags
			fs.args = fs.args[1:]
			return "", nil
		}
	}
	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return "", fs.failf("bad flag syntax: %s", s)
	}
	return name, nil
}

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the Set
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help or -h were set but not defined.
func (fs *Set) Parse(arguments []string) error {
	fs.args = arguments
	for {
		seen, err := fs.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}
		return err
	}
	return nil
}

// Sequence returns the parsed flag sequence.
func (fs *Set) Sequence() []*flag.Flag {
	return fs.actual
}

// NewSet returns a new, empty flag set with the specified name.
func NewSet(name string) *Set {
	return &Set{name: name, errorHandling: flag.ContinueOnError}
}
