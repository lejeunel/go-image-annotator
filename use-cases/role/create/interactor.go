package create

import (
	"context"
	"fmt"
	"slices"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	v "github.com/lejeunel/go-image-annotator/modules/string-validator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	Repo
	v.Validator
	Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "creating role"
	if err := i.Auth.CreateRole(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.validate(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	existingMethods := i.Auth.ListMethods()
	for _, m := range r.Methods {
		if !slices.Contains(existingMethods, m) {
			out.Error(
				fmt.Errorf(
					"%v: checking whether method %v is allowed: %w",
					errCtx,
					m,
					e.ErrValidation,
				),
			)
			return
		}
	}

	role, err := i.create(r)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessCreateRole(*role)
}

func (i *Interactor) create(r Request) (*rl.Role, error) {
	role := rl.Role{
		Id:          rl.NewRoleId(),
		Name:        r.Name,
		Description: r.Description,
		Methods:     r.Methods,
	}
	if err := i.Repo.Create(role); err != nil {
		return nil, err
	}
	return &role, nil
}

func (i *Interactor) validate(name string) error {
	if err := i.Validator.Validate(name); err != nil {
		return fmt.Errorf("checking role name %v: %w", name, err)
	}
	if err := i.isDuplicate(name); err != nil {
		return err
	}
	return nil
}

func (i *Interactor) isDuplicate(name string) error {
	errBaseMsg := fmt.Sprintf("checking for duplicate role with name %v", name)
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

func New(r Repo, a Auth, opts ...Option) Interactor {
	i := &Interactor{
		Repo: r, Validator: v.NewNameValidator(),
		Auth: a,
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
