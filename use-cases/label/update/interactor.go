package update

import (
	"context"
	"fmt"

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
type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func NewInteractor(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r,
		logger: logging.NewNoOpLogger(),
		auth:   auth.PassThroughAuth{}}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {

	if err := i.auth.UpdateLabel(ctx); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.ensureNameExists(r.Name); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.repo.Update(Model{Name: r.Name, NewDescription: r.NewDescription}); err != nil {
		i.handleError(err, out)
		return
	}

	out.Success(Response{Name: r.Name, Description: r.NewDescription})
}
func (i *Interactor) ensureNameExists(name string) error {
	exists, err := i.repo.Exists(name)
	errCtx := fmt.Errorf("checking that label %v exists", name)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, e.ErrInternal)
	}
	if !exists {
		return fmt.Errorf("%w: %w", errCtx, e.ErrNotFound)
	}
	return nil
}
func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "updating label"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}
