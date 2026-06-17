package list

import (
	"context"
	"fmt"

	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Interactor struct {
	repo Repo
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "listing groups"
	if err := pagination.Validate(r.Page, r.PageSize); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	found, err := i.repo.List(r)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	count, err := i.repo.Count()
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	response := Response{Pagination: pagination.New(int64(r.Page), r.PageSize, *count)}
	response.Groups = found
	out.Success(response)
}

type Option func(*Interactor)

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
