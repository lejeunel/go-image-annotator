package renew_token

import (
	"context"
	"fmt"

	"log/slog"

	g "github.com/lejeunel/go-image-annotator/app/token-generator"
	usr "github.com/lejeunel/go-image-annotator/entities/user"
	"github.com/lejeunel/go-image-annotator/shared/auth"
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
	if err := i.auth.RenewToken(ctx, r.Id); err != nil {
		i.handleError(err, out)
		return

	}
	token, err := i.tokenGenerator.Generate()
	if err != nil {
		i.handleError(err, out)
		return
	}

	user := usr.NewUser(r.Id, usr.WithHashedPersonalAccessToken(token.Hash))
	if err := i.repo.Create(user); err != nil {
		i.handleError(err, out)
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

func NewInteractor(r Repo, g TokenGenerator, opts ...Option) *Interactor {
	i := &Interactor{repo: r,
		logger:         logging.NewNoOpLogger(),
		auth:           auth.PassThroughAuth{},
		tokenGenerator: g,
	}

	for _, opt := range opts {
		opt(i)
	}
	return i
}
