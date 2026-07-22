package delete

import (
	"fmt"

	"context"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	repo Repo
	auth Auth
}

func (i *Interactor) Execute(ctx context.Context, id string, out OutputPort) {
	errCtx := "deleting user"
	if err := i.auth.DeleteUser(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if err := i.exists(id); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	currentUser := u.IdentityFromContext(ctx)
	if currentUser == nil {
		out.Error(fmt.Errorf("%v: fetching user from context: %w", errCtx, e.ErrInternal))
		return
	}

	if currentUser.Id == id {
		out.Error(fmt.Errorf("%v: attempting to delete user %v while logged-in as %v: %w", errCtx, id, currentUser.Id, e.ErrForbiddenOp))
		return
	}

	if err := i.repo.Delete(id); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.SuccessDeleteUser(id)
}
func (i *Interactor) exists(name string) error {
	errCtx := fmt.Errorf("checking whether user with id %v exists", name)
	exists, err := i.repo.Exists(name)
	if err != nil {
		return fmt.Errorf("%w: %v: %w", errCtx, err, e.ErrInternal)
	}
	if !exists {
		return fmt.Errorf("%w: %v: %w", errCtx, err, e.ErrNotFound)
	}
	return nil
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r,
		auth: auth.NewVoidAuth()}
	for _, opt := range opts {
		opt(i)
	}
	return *i

}
