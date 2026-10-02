package find

import (
	"context"
	"fmt"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
)

type Interactor struct {
	UserStore
	Auth
}

func (i *Interactor) Execute(ctx context.Context, id u.UserId, out OutputPort) {
	errCtx := "fetching user"
	if err := i.Auth.FindUser(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	found, err := i.UserStore.Find(id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessFindUser(*found)
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(r UserStore, opts ...Option) Interactor {
	i := &Interactor{
		UserStore: r,
		Auth:      auth.NewVoidAuth(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
