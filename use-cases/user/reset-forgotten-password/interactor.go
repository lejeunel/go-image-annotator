package reset_forgotten_password

import (
	"context"
	"fmt"

	"github.com/jonboulle/clockwork"
	tk "github.com/lejeunel/go-image-annotator/modules/token"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type PasswordValidator interface {
	Validate(password string) error
}

type Interactor struct {
	repo              Repo
	tokenHasher       tk.TokenHasher
	passwordValidator PasswordValidator
	clock             clockwork.Clock
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "resetting forgotten password"

	if r.FirstPassword != r.SecondPassword {
		out.Error(fmt.Errorf("%v: checking for matching passwords: %w", errCtx, e.ErrPasswordMismatch))
		return
	}

	if err := i.passwordValidator.Validate(r.FirstPassword); err != nil {
		out.Error(fmt.Errorf("%v: checking for password validity: %w", errCtx, e.ErrInvalidPassword))
		return
	}

	state, err := i.repo.FindResetPasswordState(i.tokenHasher.Hash(r.Token))
	if err != nil {
		out.Error(fmt.Errorf("%v: finding by hash: %w", errCtx, err))
		return
	}

	if state.ExpiresAt != nil {
		if state.ExpiresAt.Before(i.clock.Now()) {
			out.Error(fmt.Errorf("%v: checking for token expiration: %w", errCtx, e.ErrExpiredToken))
			return
		}
	}

	if err := i.repo.UpdatePassword(state.Id, i.tokenHasher.Hash(r.FirstPassword)); err != nil {
		out.Error(fmt.Errorf("%v: updating password: %v, %w", errCtx, err, e.ErrInternal))
		return
	}

	out.Success()
}

type Option func(*Interactor)

func WithClock(c clockwork.Clock) Option {
	return func(i *Interactor) {
		i.clock = c
	}
}

func New(r Repo, tokenHasher tk.TokenHasher, passwordValidator PasswordValidator,
	opts ...Option) Interactor {
	i := &Interactor{repo: r,
		tokenHasher:       tokenHasher,
		passwordValidator: passwordValidator,
		clock:             clockwork.NewRealClock(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
