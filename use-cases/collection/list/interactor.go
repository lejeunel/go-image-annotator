package list

import (
	"context"
	"fmt"

	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Interactor struct {
	Repo
	DefaultPageSize int
	MaxPageSize     int
}

func (i Interactor) Execute(ctx context.Context, r pa.PaginationParams, out OutputPort) {
	errCtx := "listing collections"
	r.Sanitize(i.DefaultPageSize, i.MaxPageSize)

	found, err := i.Repo.List(r)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	count, err := i.Repo.Count()
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

func New(r Repo, dps int, mps int, opts ...Option) Interactor {
	i := &Interactor{Repo: r, DefaultPageSize: dps, MaxPageSize: mps}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
