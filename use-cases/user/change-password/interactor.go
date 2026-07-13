package change_password

import (
	"context"
	"fmt"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	pw "github.com/lejeunel/go-image-annotator/modules/password-validator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type TokenHasher interface {
	Hash(token string) []byte
}

type Interactor struct {
	repo              Repo
	tokenHasher       TokenHasher
	passwordValidator pw.PasswordValidator
	auth              Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "changing password"

	if err := i.auth.ChangePassword(ctx, r.Id); err != nil {
		out.Error(fmt.Errorf("%v: checking authorization: %w", errCtx, err))
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

	if err := i.repo.UpdatePassword(r.Id, i.tokenHasher.Hash(r.FirstPassword)); err != nil {
		out.Error(fmt.Errorf("%v: updating password: %v, %w", errCtx, err, e.ErrInternal))
		return
	}

	out.Success()
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo, tokenHasher TokenHasher, passwordValidator pw.PasswordValidator,
	opts ...Option) Interactor {
	i := &Interactor{repo: r,
		tokenHasher:       tokenHasher,
		passwordValidator: passwordValidator,
		auth:              auth.NewVoidAuth(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
