package assign_group

import (
	"context"
	"fmt"
	"slices"

	"log/slog"

	"github.com/lejeunel/go-image-annotator/shared/auth"
	"github.com/lejeunel/go-image-annotator/shared/logging"
)

type Interactor struct {
	repo   Repo
	logger *slog.Logger
	auth   Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	if err := i.auth.AssignUserToGroup(ctx, r.Id, r.Group); err != nil {
		i.handleError(err, out)
		return

	}
	if err := i.repo.UserExists(r.Id); err != nil {
		i.handleError(err, out)
		return
	}
	if err := i.repo.GroupExists(r.Id); err != nil {
		i.handleError(err, out)
		return
	}
	user, err := i.repo.Find(r.Id)
	if err != nil {
		i.handleError(err, out)
		return
	}
	if slices.Contains(user.Groups, r.Group) {
		out.Success(Response{Id: r.Id, Groups: user.Groups})
		return
	}
	if err := i.repo.AssignToGroup(r.Id, r.Group); err != nil {
		i.handleError(err, out)
		return
	}
	out.Success(Response{Id: r.Id, Groups: append(user.Groups, r.Group)})
}
func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "creating user"
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

func NewInteractor(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r,
		logger: logging.NewNoOpLogger(),
		auth:   auth.PassThroughAuth{},
	}

	for _, opt := range opts {
		opt(i)
	}
	return i
}
