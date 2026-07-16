package list

import (
	"context"
	"fmt"

	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Interactor struct {
	repo Repo
}

func (i Interactor) Execute(ctx context.Context, r pa.PaginationParams, out OutputPort) {
	errCtx := "listing collections"
	if err := pa.Validate(r.Page, r.PageSize); err != nil {
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

	response := Response{Pagination: pa.New(int64(r.Page), r.PageSize, *count)}
	for _, f := range found {
		response.Collections = append(response.Collections, *f)
	}
	out.SuccessListCollections(response)
}

type Option func(*Interactor)

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
