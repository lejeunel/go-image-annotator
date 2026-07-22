package update_group

import (
	"context"
	"fmt"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
)

type Interactor struct {
	userRepo  UserRepo
	groupRepo GroupRepo
	auth      Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "assigning group to user"
	if err := i.auth.UpdateGroups(ctx); err != nil {
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

	out.SuccessUpdateGroups(Response{Id: r.Id, Groups: r.Groups})
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
