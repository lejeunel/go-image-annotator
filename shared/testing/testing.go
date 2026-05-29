package testing

import (
	"errors"
	"reflect"
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type TestingErrPresenter struct {
	GotDuplicationErr bool
	GotValidationErr  bool
	GotInternalErr    bool
	GotNotFoundErr    bool
	GotDependencyErr  bool
	GotErr            error
}

func (p *TestingErrPresenter) Error(err error) {
	p.GotErr = err
	switch {
	case errors.Is(err, e.ErrDuplicate):
		p.GotDuplicationErr = true
	case errors.Is(err, e.ErrValidation):
		p.GotValidationErr = true
	case errors.Is(err, e.ErrNotFound):
		p.GotNotFoundErr = true
	case errors.Is(err, e.ErrDependency):
		p.GotDependencyErr = true

	default:
		p.GotInternalErr = true
	}
}

func AssertEqual(t *testing.T, prefixMsg string, got any, want any) {
	if got != want {
		t.Fatalf("%v: got %v, want %v",
			prefixMsg, got, want)
	}
}

func AssertContains[T comparable](t *testing.T, prefixMsg string, got []T, want T) {
	for _, g := range got {
		if g == want {
			return
		}
	}
	t.Fatalf("%v: expected %v, to contain %v",
		prefixMsg, got, want)
}

func AssertNonNil(t *testing.T, prefixMsg string, v any) {
	rv := reflect.ValueOf(v)
	if (v == nil) || (rv.Kind() == reflect.Pointer && rv.IsNil()) {
		t.Fatalf("%v: expected non-nil variable", prefixMsg)
	}
}
