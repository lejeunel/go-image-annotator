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
	Repo
	tk.TokenHasher
	PasswordValidator
	clockwork.Clock
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "resetting forgotten password"

	if r.FirstPassword != r.SecondPassword {
		out.Error(
			fmt.Errorf("%v: checking for matching passwords: %w", errCtx, e.ErrPasswordMismatch),
		)
		return
	}

	if err := i.PasswordValidator.Validate(r.FirstPassword); err != nil {
		out.Error(
			fmt.Errorf(
				"%v: checking for password validity: %w: %w",
				errCtx,
				err,
				e.ErrInvalidPassword,
			),
		)
		return
	}

	state, err := i.Repo.FindResetPasswordState(i.TokenHasher.Hash(r.Token))
	if err != nil {
		out.Error(fmt.Errorf("%v: finding by hash: %w", errCtx, err))
		return
	}

	if state.ExpiresAt != nil {
		if state.ExpiresAt.Before(i.Clock.Now()) {
			out.Error(
				fmt.Errorf("%v: checking for token expiration: %w", errCtx, e.ErrExpiredToken),
			)
			return
		}
	}

	if err := i.Repo.UpdatePassword(state.Id, i.TokenHasher.Hash(r.FirstPassword)); err != nil {
		out.Error(fmt.Errorf("%v: updating password: %v, %w", errCtx, err, e.ErrInternal))
		return
	}
	if err := i.Repo.DeleteForgottenPasswordTokens(state.Id); err != nil {
		out.Error(fmt.Errorf("%v: deleting token: %v, %w", errCtx, err, e.ErrInternal))
		return
	}

	out.Success()
}

type Option func(*Interactor)

func WithClock(c clockwork.Clock) Option {
	return func(i *Interactor) {
		i.Clock = c
	}
}

func New(r Repo, tokenHasher tk.TokenHasher, passwordValidator PasswordValidator,
	opts ...Option,
) Interactor {
	i := &Interactor{
		Repo:              r,
		TokenHasher:       tokenHasher,
		PasswordValidator: passwordValidator,
		Clock:             clockwork.NewRealClock(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
