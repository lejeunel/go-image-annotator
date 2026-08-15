package list

import (
	"context"
	"fmt"

	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Interactor struct {
	Repo
	DefaultPageSize int
	MaxPageSize     int
}

func (i *Interactor) Execute(ctx context.Context, r pag.PaginationParams, out OutputPort) {
	errCtx := "listing labels"

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

	response := Response{
		Pagination: pag.New(int64(r.Page), r.PageSize, count),
	}
	for _, f := range found {
		response.Labels = append(response.Labels, *f)
	}
	out.SuccessListLabels(response)
}

type Option func(*Interactor)

func New(r Repo, dps int, mps int, opts ...Option) *Interactor {
	i := &Interactor{r, dps, mps}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
