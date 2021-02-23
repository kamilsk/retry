package retry

import (
	"errors"
	"reflect"
	"testing"
)

func TestError(t *testing.T) {
	if internal.Error() != string(internal) {
		t.Error("unexpected behavior")
	}

	if internal.Unwrap() != nil {
		t.Error("unexpected behavior")
	}
}

func TestUnwrap(t *testing.T) {
	root := errors.New("root")
	core := unwrap(cause{layer{root}})
	if !reflect.DeepEqual(core, root) {
		t.Error("unexpected behavior")
	}
}

// helpers

type cause struct{ error }

func (cause cause) Cause() error { return cause.error }

type layer struct{ error }

func (layer layer) Unwrap() error { return layer.error }
