package update

import (
	"context"
	"errors"
	"fmt"
	"slices"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	Repo
	Auth
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{
		Repo: r,
		Auth: auth.NewVoidAuth(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "updating role"
	if err := i.Auth.UpdateRole(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	role, err := i.Repo.Find(r.Name)
	if err != nil {
		out.Error(fmt.Errorf("%v: fetching role %v: %w", errCtx, r.Name, err))
		return
	}

	if r.NewName != r.Name {
		if err := i.ensureNameDoesNotExist(r.NewName); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	existingMethods := i.Auth.ListMethods()
	for _, m := range r.NewMethods {
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

	if err := i.Repo.Update(
		rl.UpdatableModel{
			Name:           r.Name,
			NewName:        r.NewName,
			NewDescription: r.NewDescription,
			NewMethods:     r.NewMethods,
		},
	); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessUpdateRole(
		rl.Role{Id: role.Id, Name: r.NewName, Description: r.NewDescription, Methods: r.NewMethods},
	)
}

func (i *Interactor) ensureNameDoesNotExist(name string) error {
	baseErr := fmt.Errorf("ensuring that a role with name %v does not already exist", name)
	_, err := i.Repo.Find(name)
	if errors.Is(err, e.ErrNotFound) {
		return nil
	}
	if err == nil {
		return fmt.Errorf("%w: %w", baseErr, e.ErrDuplicate)
	}
	return err
}
