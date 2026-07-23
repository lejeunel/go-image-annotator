package update

import (
	"context"
	"fmt"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	userRepo  UserRepo
	groupRepo GroupRepo
	roleRepo  RoleRepo
	auth      Auth
}

func New(ur UserRepo, gr GroupRepo, rr RoleRepo, opts ...Option) Interactor {
	i := &Interactor{userRepo: ur,
		groupRepo: gr,
		roleRepo:  rr,
		auth:      auth.NewVoidAuth(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "assigning group to user"
	if err := i.auth.UpdateUserPrivileges(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	_, err := i.userRepo.Find(r.Id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	for _, g := range r.Groups {
		exists, err := i.groupRepo.Exists(g)
		if err != nil {
			out.Error(fmt.Errorf("%v: checking whether group %v exists: %w", errCtx, g, err))
			return
		}
		if !*exists {
			out.Error(fmt.Errorf("%v: checking whether group %v exists: %w", errCtx, g, err))
			return
		}
	}

	if err := i.userRepo.SetGroups(r.Id, r.Groups); err != nil {
		out.Error(fmt.Errorf("%v: applying groups %v: %w", errCtx, r.Groups, err))
		return
	}

	for _, role := range r.Roles {
		exists, err := i.roleRepo.Exists(role)
		if err != nil {
			out.Error(fmt.Errorf("%v: checking whether role %v exists: %w", errCtx, role, err))
			return
		}
		if !*exists {
			out.Error(fmt.Errorf("%v: checking whether role %v exists: %w", errCtx, role, e.ErrNotFound))
			return
		}
	}
	if err := i.userRepo.SetRoles(r.Id, r.Roles); err != nil {
		out.Error(fmt.Errorf("%v: applying roles %v: %w", errCtx, r.Roles, err))
		return
	}

	out.SuccessUpdate(Response{Id: r.Id, Groups: r.Groups})
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}
