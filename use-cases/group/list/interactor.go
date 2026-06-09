package list

import (
	"context"
	"fmt"

	"github.com/lejeunel/go-image-annotator/shared/logging"
	"github.com/lejeunel/go-image-annotator/shared/pagination"

	"log/slog"
)

type Interactor struct {
	repo   Repo
	logger *slog.Logger
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	if err := pagination.Validate(r.Page, r.PageSize); err != nil {
		i.handleError(err, out)
		return
	}

	found, err := i.repo.List(r)
	if err != nil {
		i.handleError(err, out)
		return
	}

	count, err := i.repo.Count()
	if err != nil {
		i.handleError(err, out)
		return
	}

	response := Response{Pagination: pagination.New(int64(r.Page), r.PageSize, *count)}
	response.Groups = found
	out.Success(response)
}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "listing images"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

type Option func(*Interactor)

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		logger: logging.NewNoOpLogger(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
