package unassign_role

import (
	"context"
	"fmt"
	"slices"

	"github.com/lejeunel/go-image-annotator/shared/auth"
)

type Interactor struct {
	repo Repo
	auth Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "creating user"
	if err := i.auth.UnAssignRoleFromUser(ctx, r.Id, r.Role); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return

	}
	user, err := i.repo.Find(r.Id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if !slices.Contains(user.Roles, r.Role) {
		out.Success(Response{Id: r.Id, Roles: user.Roles})
		return
	}
	if err := i.repo.UnAssignRole(r.Id, r.Role); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
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

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		auth: auth.PassThroughAuth{},
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
