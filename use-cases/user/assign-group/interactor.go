package assign_group

import (
	"context"
	"fmt"
	"slices"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
)

type Interactor struct {
	userRepo  UserRepo
	groupRepo GroupRepo
	auth      Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "assigning group to user"
	if err := i.auth.AssignUserToGroup(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	exists, err := i.groupRepo.Exists(r.Group)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if !*exists {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	user, err := i.userRepo.Find(r.Id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if slices.Contains(user.Groups, r.Group) {
		out.Success(Response{Id: r.Id, Groups: user.Groups})
		return
	}
	if err := i.userRepo.AssignToGroup(r.Id, r.Group); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.Success(Response{Id: r.Id, Groups: append(user.Groups, r.Group)})
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(ur UserRepo, gr GroupRepo, opts ...Option) Interactor {
	i := &Interactor{userRepo: ur,
		groupRepo: gr,
		auth:      auth.NewVoidAuth(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
