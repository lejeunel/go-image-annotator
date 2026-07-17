package list

import (
	"context"
	"fmt"

	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Interactor struct {
	repo Repo
}

func (i *Interactor) Execute(ctx context.Context, r pag.PaginationParams, out OutputPort) {
	errCtx := "listing labels"
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

	response := Response{
		Pagination: pag.New(int64(r.Page), r.PageSize, count),
	}
	for _, f := range found {
		response.Labels = append(response.Labels, *f)
	}
	out.SuccessListLabels(response)

}

type Option func(*Interactor)

func New(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
