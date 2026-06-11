package renew_token

import (
	"context"
	"fmt"

	"log/slog"

	g "github.com/lejeunel/go-image-annotator/app/token-generator"
	"github.com/lejeunel/go-image-annotator/shared/auth"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/logging"
)

type TokenGenerator interface {
	Generate() (*g.TokenPair, error)
}

type Interactor struct {
	repo           Repo
	logger         *slog.Logger
	tokenGenerator TokenGenerator

	auth Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := fmt.Errorf("setting api access token to user %v", r.Id)
	if err := i.auth.RenewToken(ctx, r.Id); err != nil {
		i.handleError(err, out)
		return

	}
	exists, err := i.repo.Exists(r.Id)
	if err != nil {
		out.Error(fmt.Errorf("%w: checking user exists: %w", errCtx, err))
		return
	}
	if !exists {
		out.Error(fmt.Errorf("%w: checking user exists: %w", errCtx, e.ErrNotFound))
		return
	}

	token, err := i.tokenGenerator.Generate()
	if err != nil {
		out.Error(fmt.Errorf("%w: generating token: %w", errCtx, err))
		return
	}

	if err := i.repo.SetAccessTokenHash(r.Id, token.Hash); err != nil {
		out.Error(fmt.Errorf("%w: setting token hash: %w", errCtx, err))
		return
	}
	out.Success(Response{Id: r.Id, PersonalAccessToken: token.Token})
}
func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "renewing personal access token"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo, g TokenGenerator, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		logger:         logging.NewNoOpLogger(),
		auth:           auth.PassThroughAuth{},
		tokenGenerator: g,
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
