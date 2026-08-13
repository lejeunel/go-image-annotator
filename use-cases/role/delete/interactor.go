package delete

import (
	"context"
	"fmt"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	Repo
	Auth
}

func (i *Interactor) Execute(ctx context.Context, name string, out OutputPort) {
	errCtx := fmt.Errorf("deleting role")
	if err := i.Auth.DeleteRole(ctx); err != nil {
		out.Error(fmt.Errorf("%w: %w", errCtx, e.ErrAuthorization))
		return
	}

	if err := i.ensureExists(name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if err := i.ensureDeletable(name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.Repo.Delete(name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.SuccessDeleteRole(name)
}

func (i *Interactor) ensureDeletable(name string) error {
	errCtx := fmt.Errorf("ensuring role with name %v is not assigned to any user", name)
	isAssigned, err := i.Repo.IsAssigned(name)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, e.ErrInternal)
	}
	if *isAssigned {
		return fmt.Errorf("%w: %w", errCtx, e.ErrDependency)
	}
	return nil
}

func (i *Interactor) ensureExists(name string) error {
	errCtx := fmt.Errorf("checking whether role with name %v exists", name)
	exists, err := i.Repo.Exists(name)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, e.ErrInternal)
	}
	if !*exists {
		return fmt.Errorf("%w: %w", errCtx, e.ErrNotFound)
	}
	return nil
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
