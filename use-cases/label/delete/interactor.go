package delete

import (
	"fmt"

	"context"
	"log/slog"

	"github.com/lejeunel/go-image-annotator/shared/auth"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/logging"
)

type Interactor struct {
	repo   Repo
	logger *slog.Logger
	auth   Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	if err := i.auth.DeleteLabel(ctx); err != nil {
		i.handleError(err, out)
		return
	}
	if err := i.isUsed(r.Name); err != nil {
		i.handleError(err, out)
		return
	}
	if err := i.exists(r.Name); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.repo.Delete(r.Name); err != nil {
		i.handleError(err, out)
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

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "creating label"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func NewInteractor(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r, logger: logging.NewNoOpLogger(),
		auth: auth.PassThroughAuth{}}
	for _, opt := range opts {
		opt(i)
	}
	return i

}
