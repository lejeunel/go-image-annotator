package set_admin

import (
	"context"
	"fmt"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
)

type Interactor struct {
	repo Repo
	auth Auth
}

func (i *Interactor) Execute(ctx context.Context, id u.UserId, value bool, out OutputPort) {
	errCtx := fmt.Errorf("settings admin rights to %v for user %v", id, value)
	if err := i.auth.SetAdminRights(ctx); err != nil {
		out.Error(fmt.Errorf("%w: %w", errCtx, err))
		return
	}
	if err := i.repo.SetAdmin(id, value); err != nil {
		out.Error(fmt.Errorf("%w: %w", errCtx, err))
		return
	}
	out.Success(Response{Id: id, IsAdmin: value})
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		auth: auth.NewVoidAuth(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
