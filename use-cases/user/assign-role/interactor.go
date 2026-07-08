package assign_role

import (
	"context"
	"fmt"
	"slices"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
)

type Interactor struct {
	repo UserRepo
	auth Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "creating user"
	if err := i.auth.AssignRoleToUser(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return

	}
	user, err := i.repo.Find(r.Id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if slices.Contains(user.Roles, r.Role) {
		out.Success(Response{Id: r.Id, Roles: user.Roles})
		return
	}
	if err := i.repo.AssignRole(r.Id, r.Role); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.Success(Response{Id: r.Id, Roles: append(user.Roles, r.Role)})
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r UserRepo, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		auth: auth.NewVoidAuth(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
