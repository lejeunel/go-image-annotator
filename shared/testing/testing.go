package testing

import (
	"errors"
	i "github.com/lejeunel/go-image-annotator/entities/identity"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type TestingErrPresenter struct {
	GotDuplicationErr bool
	GotValidationErr  bool
	GotInternalErr    bool
	GotNotFoundErr    bool
	GotDependencyErr  bool
	GotErr            error
	GotAuthErr        bool
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
	case errors.Is(err, e.ErrAuth):
		p.GotAuthErr = true

	default:
		p.GotInternalErr = true
	}
}

type FakeAuth struct {
	Fail bool
}

type FakeProvider struct {
}

func (p FakeProvider) Provide() (*i.Identity, error) {
	return &i.Identity{}, nil
}
