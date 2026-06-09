package read

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lejeunel/go-image-annotator/shared/logging"
)

type Interactor struct {
	repo   Repo
	logger *slog.Logger
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	found, err := i.repo.FindLabel(r.Name)
	if err != nil {
		i.handleError(err, out)
		return
	}

	out.Success(Response{Name: found.Name, Description: found.Description})

}
func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "fetching label"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

type Option func(*Interactor)

func New(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r, logger: logging.NewNoOpLogger()}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
