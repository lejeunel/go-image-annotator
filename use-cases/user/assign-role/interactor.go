package assign_role

import (
	"context"
	"fmt"
	"slices"

	"log/slog"

	"github.com/lejeunel/go-image-annotator/shared/auth"
	"github.com/lejeunel/go-image-annotator/shared/logging"
)

type Interactor struct {
	repo   UserRepo
	logger *slog.Logger
	auth   Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	if err := i.auth.AssignRoleToUser(ctx, r.Id, r.Role); err != nil {
		i.handleError(err, out)
		return

	}
	user, err := i.repo.Find(r.Id)
	if err != nil {
		i.handleError(err, out)
		return
	}
	if slices.Contains(user.Roles, r.Role) {
		out.Success(Response{Id: r.Id, Roles: user.Roles})
		return
	}
	if err := i.repo.AssignRole(r.Id, r.Role); err != nil {
		i.handleError(err, out)
		return
	}
	out.Success(Response{Id: r.Id, Roles: append(user.Roles, r.Role)})
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

func NewInteractor(r UserRepo, opts ...Option) *Interactor {
	i := &Interactor{repo: r,
		logger: logging.NewNoOpLogger(),
		auth:   auth.PassThroughAuth{},
	}

	for _, opt := range opts {
		opt(i)
	}
	return i
}
