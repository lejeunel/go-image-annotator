package unassign_group

import (
	"context"
	"fmt"

	"github.com/lejeunel/go-image-annotator/modules/auth"
)

type Interactor struct {
	userRepo  UserRepo
	groupRepo GroupRepo
	auth      Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "un-assigning user from group"
	if err := i.auth.UnAssignUserFromGroup(ctx); err != nil {
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
	if err := i.userRepo.UnAssignFromGroup(r.Id, r.Group); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	newGroups := []string{}
	for _, g := range user.Groups {
		if g == r.Group {
			continue
		}
		newGroups = append(newGroups, g)
	}
	out.Success(Response{Id: r.Id, Groups: newGroups})
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
