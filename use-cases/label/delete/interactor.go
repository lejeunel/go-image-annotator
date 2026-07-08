package delete

import (
	"context"
	"fmt"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	repo Repo
	auth Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "deleting label"
	if err := i.auth.DeleteLabel(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if err := i.isUsed(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if err := i.exists(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.repo.Delete(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.Success()
}
func (i *Interactor) exists(name string) error {
	errCtx := fmt.Errorf("checking whether label with name %v exists", name)
	exists, err := i.repo.Exists(name)
	if err != nil {
		return fmt.Errorf("%w: %v: %w", errCtx, err, e.ErrInternal)
	}
	if !exists {
		return fmt.Errorf("%w: %v: %w", errCtx, err, e.ErrNotFound)
	}
	return nil
}

func (i *Interactor) isUsed(name string) error {
	errCtx := fmt.Errorf("checking whether label with name %v is used", name)
	isUsed, err := i.repo.IsUsed(name)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, e.ErrInternal)
	}
	if *isUsed {
		return fmt.Errorf("%w: %w", errCtx, e.ErrDependency)
	}
	return nil

}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r,
		auth: auth.NewVoidAuth()}
	for _, opt := range opts {
		opt(i)
	}
	return i

}
