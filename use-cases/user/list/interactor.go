package list

import (
	"context"
	"fmt"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Interactor struct {
	Repo
	Auth
}

func (i *Interactor) Execute(ctx context.Context, r pag.PaginationParams, out OutputPort) {
	errCtx := "listing users"

	if err := i.Auth.ListUsers(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
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
		Users:      found,
	}
	out.SuccessListUsers(response)
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{
		Repo: r,
		Auth: auth.NewVoidAuth(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
