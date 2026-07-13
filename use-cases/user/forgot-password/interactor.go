package forgot_password

import (
	"context"
	"fmt"
	"github.com/jonboulle/clockwork"
	"time"

	tk "github.com/lejeunel/go-image-annotator/entities/token"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type TokenGenerator interface {
	Generate() (*tk.Token, error)
}

type Interactor struct {
	repo           Repo
	expiresMinutes int
	tokenGenerator TokenGenerator
	auth           Auth
	clock          clockwork.Clock
}

func (i *Interactor) Execute(ctx context.Context, userId string, out OutputPort) {
	errCtx := "requesting forgotten password token"
	if err := i.auth.RequestForgottenPasswordToken(ctx); err != nil {
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

	if err := i.repo.DeleteForgottenPasswordTokens(userId); err != nil {
		out.Error(fmt.Errorf("%v: deleting previous tokens: %w", errCtx, e.ErrInternal))
		return
	}

	token, err := i.tokenGenerator.Generate()
	if err != nil {
		out.Error(fmt.Errorf("%v: generating token: %w", errCtx, err))
		return
	}

	expiresAt := i.clock.Now().Add(time.Minute * time.Duration(i.expiresMinutes))
	if err := i.repo.AddForgottenPasswordState(token.Hash, userId, expiresAt); err != nil {
		out.Error(fmt.Errorf("%v: storing token: %w", errCtx, err))
		return
	}
	out.Success(Response{Id: userId, Email: userId,
		PasswordResetToken: token.Value})
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func WithClock(c clockwork.Clock) Option {
	return func(i *Interactor) {
		i.clock = c
	}
}

func New(r Repo, expiresMinutes int, g TokenGenerator, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		auth:           auth.NewVoidAuth(),
		tokenGenerator: g,
		expiresMinutes: expiresMinutes,
		clock:          clockwork.NewRealClock(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
