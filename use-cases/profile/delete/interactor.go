package delete

import (
	"context"
	"fmt"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	Repo
	Auth
}

func (i *Interactor) Execute(ctx context.Context, name string, out OutputPort) {
	errCtx := "deleting profile"

	profile, err := i.Repo.Find(name)
	if err != nil {
		out.Error(fmt.Errorf("%v: fetching collection: %w", errCtx, err))
		return
	}

	if err := i.Auth.DeleteProfile(ctx, profile.Group); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	isUsed, err := i.Repo.IsUsed(name)
	if err != nil {
		out.Error(fmt.Errorf("%v: checking whether profile is used: %w", errCtx, err))
		return
	}
	if *isUsed {
		out.Error(fmt.Errorf("%v: checking whether profile is used: %w", errCtx, e.ErrDependency))
		return
	}
	if err := i.Repo.Delete(name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.SuccessDeleteProfile(name)
}

// func (i *Interactor) exists(name string) error {
// 	errCtx := fmt.Errorf("checking whether label with name %v exists", name)
// 	exists, err := i.Repo.Exists(name)
// 	if err != nil {
// 		return fmt.Errorf("%w: %v: %w", errCtx, err, e.ErrInternal)
// 	}
// 	if !exists {
// 		return fmt.Errorf("%w: %v: %w", errCtx, err, e.ErrNotFound)
// 	}
// 	return nil
// }

// func (i *Interactor) isUsed(name string) error {
// 	errCtx := fmt.Errorf("checking whether label with name %v is used", name)
// 	isUsed, err := i.Repo.IsUsed(name)
// 	if err != nil {
// 		return fmt.Errorf("%w: %w", errCtx, e.ErrInternal)
// 	}
// 	if *isUsed {
// 		return fmt.Errorf("%w: %w", errCtx, e.ErrDependency)
// 	}
// 	return nil
//}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{
		Repo: r,
		Auth: auth.NewVoidAuth(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
