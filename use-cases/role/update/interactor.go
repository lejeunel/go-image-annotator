package update

import (
	"context"
	"fmt"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	repo Repo
	auth Auth
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		auth: auth.NewVoidAuth()}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {

	errCtx := "updating role"
	if err := i.auth.UpdateRole(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.ensureNameExists(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return

	}

	if r.NewName != r.Name {
		if err := i.ensureNameDoesNotExist(r.NewName); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}

	}

	if err := i.repo.Update(rl.UpdatableModel{Name: r.Name, NewName: r.NewName, NewDescription: r.NewDescription}); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessUpdateRole(Response{Name: r.NewName, Description: r.NewDescription})
}

func (i *Interactor) ensureNameExists(name string) error {
	baseErr := fmt.Errorf("ensuring that role with name %v exists", name)
	exists, err := i.repo.Exists(name)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, e.ErrInternal)
	}
	if !*exists {
		return fmt.Errorf("%w: %w", baseErr, e.ErrNotFound)
	}
	return nil
}

func (i *Interactor) ensureNameDoesNotExist(name string) error {
	baseErr := fmt.Errorf("ensuring that a role with name %v does not already exist", name)
	exists, err := i.repo.Exists(name)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, e.ErrInternal)
	}
	if *exists {
		return fmt.Errorf("%w: %w", baseErr, e.ErrDuplicate)
	}
	return nil
}
