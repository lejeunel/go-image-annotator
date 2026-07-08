package find

import (
	"context"
	"fmt"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
)

type Interactor struct {
	repo Repo
	auth Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "fetching user"
	if err := i.auth.FindUser(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	found, err := i.repo.Find(r.Id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.Success(Response{Id: found.Id, Groups: found.Groups,
		Roles: found.Roles})

}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		auth: auth.NewVoidAuth()}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
