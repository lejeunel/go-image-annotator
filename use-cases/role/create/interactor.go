package create

import (
	"context"
	"fmt"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	v "github.com/lejeunel/go-image-annotator/shared/validation"
)

type Interactor struct {
	repo      Repo
	validator v.Validator
	auth      Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "creating role"
	if err := i.auth.CreateRole(ctx); err != nil {
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

	out.SuccessCreateRole(Response{Name: r.Name, Description: r.Description})
}

func (i *Interactor) create(r Request) error {
	role := rl.NewRole(rl.NewRoleId(), r.Name, rl.WithDescription(r.Description))
	if err := i.repo.Create(role); err != nil {
		return err
	}
	return nil

}

func (i *Interactor) validate(name string) error {
	if err := i.validator.Validate(name); err != nil {
		return fmt.Errorf("checking role name %v: %w", name, err)
	}
	if err := i.isDuplicate(name); err != nil {
		return err
	}
	return nil

}

func (i *Interactor) isDuplicate(name string) error {
	errBaseMsg := fmt.Sprintf("checking for duplicate role with name %v", name)
	alreadyExists, err := i.repo.Exists(name)
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
		i.validator = v
	}
}

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r, validator: v.NewNameValidator(),
		auth: auth.NewVoidAuth()}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
