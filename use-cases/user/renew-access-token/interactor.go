package renew_token

import (
	"context"
	"fmt"

	tk "github.com/lejeunel/go-image-annotator/entities/token"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type TokenGenerator interface {
	Generate() (*tk.Token, error)
}

type Interactor struct {
	repo           Repo
	tokenGenerator TokenGenerator

	auth Auth
}

func (i *Interactor) Execute(ctx context.Context, userId string, out OutputPort) {
	errCtx := "renewing personal access token"
	if err := i.auth.RenewToken(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return

	}
	exists, err := i.repo.Exists(userId)
	if err != nil {
		out.Error(fmt.Errorf("%v: checking user %v exists: %w", errCtx, userId, err))
		return
	}
	if !exists {
		out.Error(fmt.Errorf("%v: checking user %v exists: %w", errCtx, userId, e.ErrNotFound))
		return
	}

	token, err := i.tokenGenerator.Generate()
	if err != nil {
		out.Error(fmt.Errorf("%v: generating token: %w", errCtx, err))
		return
	}

	if err := i.repo.SetAccessTokenHash(userId, token.Hash); err != nil {
		out.Error(fmt.Errorf("%v: setting token hash: %w", errCtx, err))
		return
	}
	out.Success(Response{Id: userId, PersonalAccessToken: token.Value})
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo, g TokenGenerator, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		auth:           auth.NewVoidAuth(),
		tokenGenerator: g,
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
