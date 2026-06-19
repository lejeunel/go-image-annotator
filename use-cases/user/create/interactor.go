package create

import (
	"context"
	"fmt"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
	"github.com/lejeunel/go-image-annotator/modules/auth"
	au "github.com/lejeunel/go-image-annotator/modules/authentifier"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type TokenGenerator interface {
	Generate() (*au.Pair, error)
}

type PasswordGenerator interface {
	Generate() (*au.Pair, error)
	Hash(string) []byte
}

type Interactor struct {
	repo              Repo
	tokenGenerator    TokenGenerator
	passwordGenerator PasswordGenerator

	auth Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "creating user"
	if err := i.auth.CreateUser(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return

	}
	if err := i.checkDuplicate(r.Id); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	token, err := i.tokenGenerator.Generate()
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	var passwordHash []byte
	if r.Password != nil {
		passwordHash = i.passwordGenerator.Hash(*r.Password)

	} else {
		passwordPair, err := i.passwordGenerator.Generate()
		passwordHash = passwordPair.Hash
		if err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}
	user := usr.NewUser(r.Id, usr.WithHashedPersonalAccessToken(token.Hash),
		usr.WithHashedPassword(passwordHash),
		usr.WithAdmin(r.IsAdmin))
	if err := i.repo.Create(user); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.Success(Response{
		Id:      r.Id,
		IsAdmin: r.IsAdmin})
}
func (i *Interactor) checkDuplicate(id string) error {
	errBaseMsg := "checking for duplicate user with id %v: %w"
	alreadyExists, err := i.repo.Exists(id)
	if err != nil {
		return fmt.Errorf(errBaseMsg, id, e.ErrInternal)
	}
	if alreadyExists {
		return fmt.Errorf(errBaseMsg, id, e.ErrDuplicate)
	}
	return nil
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo,
	tg TokenGenerator,
	pg PasswordGenerator, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		auth:              auth.NewVoidAuth(),
		tokenGenerator:    tg,
		passwordGenerator: pg,
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
