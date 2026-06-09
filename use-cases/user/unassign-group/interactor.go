package unassign_group

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/lejeunel/go-image-annotator/shared/auth"
	"github.com/lejeunel/go-image-annotator/shared/logging"
)

type Interactor struct {
	userRepo  UserRepo
	groupRepo GroupRepo
	logger    *slog.Logger
	auth      Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	if err := i.auth.UnAssignUserFromGroup(ctx, r.Id, r.Group); err != nil {
		i.handleError(err, out)
		return
	}
	exists, err := i.groupRepo.Exists(r.Group)
	if err != nil {
		i.handleError(err, out)
		return
	}
	if !*exists {
		i.handleError(fmt.Errorf("checking for existence of group %v", r.Group), out)
		return
	}
	user, err := i.userRepo.Find(r.Id)
	if err != nil {
		i.handleError(err, out)
		return
	}
	if err := i.userRepo.UnAssignFromGroup(r.Id, r.Group); err != nil {
		i.handleError(err, out)
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

func NewInteractor(ur UserRepo, gr GroupRepo, opts ...Option) *Interactor {
	i := &Interactor{userRepo: ur,
		groupRepo: gr,
		logger:    logging.NewNoOpLogger(),
		auth:      auth.PassThroughAuth{},
	}

	for _, opt := range opts {
		opt(i)
	}
	return i
}
