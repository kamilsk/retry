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
	cause := errors.New("root")
	core := unwrap(causer{layer{cause}})
	if !reflect.DeepEqual(core, cause) {
		t.Error("unexpected behavior")
	}
}

// helpers

type causer struct{ error }

func (causer causer) Cause() error { return causer.error }

type layer struct{ error }

func (layer layer) Unwrap() error { return layer.error }
