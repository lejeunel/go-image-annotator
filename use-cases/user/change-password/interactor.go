package change_password

import (
	"context"
	"fmt"

	pw "github.com/lejeunel/go-image-annotator/modules/password-validator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type TokenVerifier interface {
	Hash(token string) []byte
	Verify(string, []byte) bool
}

type Interactor struct {
	repo              Repo
	tokenVerifier     TokenVerifier
	passwordValidator pw.PasswordValidator
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "changing password"

	user, err := i.repo.Find(r.Id)
	if err != nil {
		out.Error(fmt.Errorf("%v: retrieving user info: %w", errCtx, err))
		return
	}

	if ok := i.tokenVerifier.Verify(r.CurrentPassword, user.HashPassword); !ok {
		out.Error(fmt.Errorf("%v: verifying current password: %w", errCtx, e.ErrInvalidPassword))
		return
	}

	if r.FirstPassword != r.SecondPassword {
		out.Error(fmt.Errorf("%v: checking for matching passwords: %w", errCtx, e.ErrPasswordMismatch))
		return
	}

	if err := i.passwordValidator.Validate(r.FirstPassword); err != nil {
		out.Error(fmt.Errorf("%v: checking for password validity: %w: %w", errCtx, err, e.ErrInvalidPassword))
		return
	}

	if err := i.repo.UpdatePassword(r.Id, i.tokenVerifier.Hash(r.FirstPassword)); err != nil {
		out.Error(fmt.Errorf("%v: updating password: %v, %w", errCtx, err, e.ErrInternal))
		return
	}

	out.Success()
}

type Option func(*Interactor)

func New(r Repo, tokenHasher TokenVerifier, passwordValidator pw.PasswordValidator,
	opts ...Option) Interactor {
	i := &Interactor{repo: r,
		tokenVerifier:     tokenHasher,
		passwordValidator: passwordValidator,
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
