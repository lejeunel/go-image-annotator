package create

import (
	"context"
	"fmt"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
	tok "github.com/lejeunel/go-image-annotator/modules/token"
	"github.com/lejeunel/go-image-annotator/shared/auth"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type TokenGenerator interface {
	Generate() (*tok.TokenPair, error)
}

type Interactor struct {
	repo           Repo
	tokenGenerator TokenGenerator

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

	user := usr.NewUser(r.Id, usr.WithHashedPersonalAccessToken(token.Hash),
		usr.WithAdmin(r.IsAdmin))
	if err := i.repo.Create(user); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.Success(Response{
		Id:                  r.Id,
		PersonalAccessToken: token.Token,
		IsAdmin:             r.IsAdmin})
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

func New(r Repo, g TokenGenerator, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		auth:           auth.PassThroughAuth{},
		tokenGenerator: g,
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
