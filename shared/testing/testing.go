package testing

import (
	"errors"

	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
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
