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
	if err := i.auth.DeleteUser(ctx); err != nil {
		i.handleError(err, out)
		return
	}
	if err := i.exists(r.Id); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.repo.Delete(r.Id); err != nil {
		i.handleError(err, out)
		return
	}
	out.Success()
}
func (i *Interactor) exists(name string) error {
	errCtx := fmt.Errorf("checking whether user with id %v exists", name)
	exists, err := i.repo.Exists(name)
	if err != nil {
		return fmt.Errorf("%w: %v: %w", errCtx, err, e.ErrInternal)
	}
	if !exists {
		return fmt.Errorf("%w: %v: %w", errCtx, err, e.ErrNotFound)
	}
	return nil
}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "creating user"
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

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r, logger: logging.NewNoOpLogger(),
		auth: auth.PassThroughAuth{}}
	for _, opt := range opts {
		opt(i)
	}
	return *i

}
