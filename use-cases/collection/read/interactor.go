package read

import (
	"context"
	"fmt"
	"log/slog"

	auth "github.com/lejeunel/go-image-annotator/shared/auth"
	"github.com/lejeunel/go-image-annotator/shared/logging"
)

type Interactor struct {
	repo   Repo
	logger *slog.Logger
	auth   Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	if err := i.auth.ReadCollection(ctx); err != nil {
		i.handleError(err, out)
		return
	}

	found, err := i.repo.FindCollectionByName(r.Name)
	if err != nil {
		i.handleError(err, out)
		return
	}

	out.Success(Response{Name: found.Name, Description: found.Description})

}
func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "fetching collection"
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

func NewInteractor(r Repo, opts ...Option) Interactor {
	i := &Interactor{
		repo:   r,
		logger: logging.NewNoOpLogger(),
		auth:   auth.PassThroughAuth{},
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
