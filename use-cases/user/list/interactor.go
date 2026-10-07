package list

import (
	"context"
	"fmt"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Interactor struct {
	Store
	Auth
}

func (i *Interactor) Execute(ctx context.Context, r pag.PaginationParams, out OutputPort) {
	errCtx := "listing users"

	if err := i.Auth.ListUsers(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	users, pagination, err := i.Store.List(r)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessListUsers(Response{Users: users, Pagination: *pagination})
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(r Store, opts ...Option) Interactor {
	i := &Interactor{
		Store: r,
		Auth:  auth.NewVoidAuth(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
