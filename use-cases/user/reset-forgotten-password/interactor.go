package reset_forgotten_password

import (
	"context"
	"fmt"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type TokenHasher interface {
	Hash(token string) []byte
}

type PasswordValidator interface {
	Validate(password string) error
}

type Interactor struct {
	repo              Repo
	tokenHasher       TokenHasher
	passwordValidator PasswordValidator
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "resetting forgotten password"

	state, err := i.repo.FindForgottenPassword(i.tokenHasher.Hash(r.Token))
	if err != nil {
		out.Error(fmt.Errorf("%v: finding by hash: %w", errCtx, err))
		return
	}

	if r.FirstPassword != r.SecondPassword {
		out.Error(fmt.Errorf("%v: checking for matching passwords: %w", errCtx, e.ErrPasswordMismatch))
		return
	}

	if err := i.passwordValidator.Validate(r.FirstPassword); err != nil {
		out.Error(fmt.Errorf("%v: checking for password validity: %w", errCtx, e.ErrInvalidPassword))
		return
	}

	if err := i.repo.UpdatePassword(state.Id, i.tokenHasher.Hash(r.FirstPassword)); err != nil {
		out.Error(fmt.Errorf("%v: updating password: %v, %w", errCtx, err, e.ErrInternal))
		return
	}

	out.Success()
}

type Option func(*Interactor)

func New(r Repo, tokenHasher TokenHasher, passwordValidator PasswordValidator,
	opts ...Option) Interactor {
	i := &Interactor{repo: r,
		tokenHasher:       tokenHasher,
		passwordValidator: passwordValidator,
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
