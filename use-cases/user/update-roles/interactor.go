package update_role

import (
	"context"
	"fmt"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
)

type Interactor struct {
	userRepo UserRepo
	roleRepo RoleRepo
	auth     Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "creating user"
	if err := i.auth.UpdateRoles(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return

	}
	_, err := i.userRepo.Find(r.Id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	for _, g := range r.Roles {
		exists, err := i.roleRepo.Exists(g)
		if err != nil {
			out.Error(fmt.Errorf("%v: checking whether role %v exists: %w", errCtx, g, err))
			return
		}
		if !*exists {
			out.Error(fmt.Errorf("%v: checking whether role %v exists: %w", errCtx, g, err))
			return
		}
	}
	if err := i.userRepo.SetRoles(r.Id, r.Roles); err != nil {
		out.Error(fmt.Errorf("%v: applying roles %v: %w", errCtx, r.Roles, err))
		return
	}
	out.SuccessUpdateRoles(Response{Id: r.Id, Roles: r.Roles})
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r UserRepo, opts ...Option) Interactor {
	i := &Interactor{userRepo: r,
		auth: auth.NewVoidAuth(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
