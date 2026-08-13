package create

import (
	"context"
	"fmt"

	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	v "github.com/lejeunel/go-image-annotator/modules/string-validator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	Repo
	v.Validator
	Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "creating label"
	if err := i.Auth.CreateLabel(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return

	}
	if err := i.Validator.Validate(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if err := i.checkDuplicate(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	label := lbl.NewLabel(lbl.NewLabelId(), r.Name, lbl.WithDescription(r.Description))
	if err := i.Repo.Create(label); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.Success(Response{Name: r.Name, Description: r.Description})
}

func (i *Interactor) checkDuplicate(name string) error {
	errBaseMsg := "checking for duplicate label with name %v: %w"
	alreadyExists, err := i.Repo.Exists(name)
	if err != nil {
		return fmt.Errorf(errBaseMsg, name, e.ErrInternal)
	}
	if alreadyExists {
		return fmt.Errorf(errBaseMsg, name, e.ErrDuplicate)
	}
	return nil
}

type Option func(*Interactor)

func WithNameValidator(v v.Validator) Option {
	return func(i *Interactor) {
		i.Validator = v
	}
}

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(r Repo, opts ...Option) *Interactor {
	i := &Interactor{
		Repo: r, Validator: v.NewNameValidator(),
		Auth: auth.NewVoidAuth(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return i
}
