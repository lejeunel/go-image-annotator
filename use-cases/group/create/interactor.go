package create

import (
	"context"
	"fmt"

	g "github.com/lejeunel/go-image-annotator/entities/group"
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
	errCtx := "creating group"
	if err := i.Auth.CreateGroup(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.validate(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.create(r); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.Success(Response{Name: r.Name, Description: r.Description})
}

func (i *Interactor) create(r Request) error {
	group := g.NewGroup(g.NewGroupId(), r.Name, g.WithDescription(r.Description))
	if err := i.Repo.Create(group); err != nil {
		return err
	}
	return nil
}

func (i *Interactor) validate(name string) error {
	if err := i.Validator.Validate(name); err != nil {
		return fmt.Errorf("checking group name %v: %w", name, err)
	}
	if err := i.isDuplicate(name); err != nil {
		return err
	}
	return nil
}

func (i *Interactor) isDuplicate(name string) error {
	errBaseMsg := fmt.Sprintf("checking for duplicate group with name %v", name)
	alreadyExists, err := i.Repo.Exists(name)
	if err != nil {
		return fmt.Errorf("%v: %w", errBaseMsg, e.ErrInternal)
	}
	if *alreadyExists {
		return fmt.Errorf("%v: %w", errBaseMsg, e.ErrDuplicate)
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

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{
		Repo: r, Validator: v.NewNameValidator(),
		Auth: auth.NewVoidAuth(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
