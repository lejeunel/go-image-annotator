package unassign_role

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
	if err := i.auth.UnAssignRoleFromUser(ctx, r.Id, r.Role); err != nil {
		i.handleError(err, out)
		return

	}
	user, err := i.repo.Find(r.Id)
	if err != nil {
		i.handleError(err, out)
		return
	}
	if !slices.Contains(user.Roles, r.Role) {
		out.Success(Response{Id: r.Id, Roles: user.Roles})
		return
	}
	if err := i.repo.UnAssignRole(r.Id, r.Role); err != nil {
		i.handleError(err, out)
		return
	}

	newRoles := []string{}
	for _, g := range user.Roles {
		if g == r.Role {
			continue
		}
		newRoles = append(newRoles, g)
	}
	out.Success(Response{Id: r.Id, Roles: newRoles})
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

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		logger: logging.NewNoOpLogger(),
		auth:   auth.PassThroughAuth{},
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
