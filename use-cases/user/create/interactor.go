package create

import (
	"context"
	"fmt"

	"log/slog"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
	tok "github.com/lejeunel/go-image-annotator/modules/token"
	"github.com/lejeunel/go-image-annotator/shared/auth"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/logging"
)

type TokenGenerator interface {
	Generate() (*tok.TokenPair, error)
}

type Interactor struct {
	repo           Repo
	logger         *slog.Logger
	tokenGenerator TokenGenerator

	auth Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	if err := i.auth.CreateUser(ctx); err != nil {
		i.handleError(err, out)
		return

	}
	if err := i.checkDuplicate(r.Id); err != nil {
		i.handleError(err, out)
		return
	}

	token, err := i.tokenGenerator.Generate()
	if err != nil {
		i.handleError(err, out)
		return
	}

	user := usr.NewUser(r.Id, usr.WithHashedPersonalAccessToken(token.Hash),
		usr.WithAdmin(r.IsAdmin))
	if err := i.repo.Create(user); err != nil {
		i.handleError(err, out)
		return
	}
	out.Success(Response{
		Id:                  r.Id,
		PersonalAccessToken: token.Token,
		IsAdmin:             r.IsAdmin})
}
func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "creating user"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
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
		logger:         logging.NewNoOpLogger(),
		auth:           auth.PassThroughAuth{},
		tokenGenerator: g,
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
