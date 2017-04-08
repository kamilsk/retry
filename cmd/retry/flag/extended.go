package flag

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

// TODO move to outside in version 2.0

// The main benefits above standard flag is it supports sequence.

func NewFlagSet(name string) *FlagSet {
	return &FlagSet{name: name, errorHandling: flag.PanicOnError}
}

type FlagSet struct {
	name          string
	parsed        bool
	actual        map[string]*flag.Flag
	formal        map[string]*flag.Flag
	args          []string
	errorHandling flag.ErrorHandling
	output        io.Writer

	sequence []*flag.Flag
}

func (fs *FlagSet) Args() []string { return fs.args }

func (fs *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	fs.Var(newBoolValue(value, p), name, usage)
}

func (fs *FlagSet) Flags() []*flag.Flag {
	return fs.sequence
}

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

func (fs *FlagSet) StringVar(p *string, name string, value string, usage string) {
	fs.Var(newStringValue(value, p), name, usage)
}

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

func (fs *FlagSet) failf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(fs.out(), err)
	return err
}

func (fs *FlagSet) out() io.Writer {
	if fs.output == nil {
		return os.Stderr
	}
	return fs.output
}

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
		return false, fs.failf("flag provided but not defined: -%s", name)
	}

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
		return false, fs.failf("invalid value %q for flag -%s: %s", value, name, err)
	}

	if fs.actual == nil {
		fs.actual = make(map[string]*flag.Flag)
	}
	fs.actual[name] = f
	return true, nil
}

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
