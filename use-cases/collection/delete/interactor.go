package delete

import (
	"context"
	"fmt"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	collectionRepo CollectionRepo
	groupRepo      GroupRepo
	auth           Auth
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "deleting collection"
	if err := i.authorizeDeletion(ctx, r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.ensureExists(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if err := i.ensureDeletable(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.collectionRepo.Delete(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.Success()
}
func (i Interactor) authorizeDeletion(ctx context.Context, name string) error {
	errCtx := fmt.Errorf("checking group ownership of collection with name %v is empty", name)
	group, err := i.groupRepo.GroupOfCollection(name)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, e.ErrInternal)
	}
	if group != nil {
		if err := i.auth.DeleteCollection(ctx, *group); err != nil {
			return fmt.Errorf("%w: %w", errCtx, e.ErrAuthorization)
		}
	}
	return nil

}
func (i Interactor) ensureDeletable(name string) error {
	errCtx := fmt.Errorf("ensuring collection with name %v is empty", name)
	isPopulated, err := i.collectionRepo.IsPopulated(name)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, e.ErrInternal)
	}
	if *isPopulated {
		return fmt.Errorf("%w: %w", errCtx, e.ErrDependency)
	}
	return nil
}
func (i Interactor) ensureExists(name string) error {
	errCtx := fmt.Errorf("checking whether collection with name %v exists", name)
	exists, err := i.collectionRepo.Exists(name)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, e.ErrInternal)
	}
	if !exists {
		return fmt.Errorf("%w: %w", errCtx, e.ErrNotFound)
	}
	return nil
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(cr CollectionRepo, gr GroupRepo, opts ...Option) Interactor {
	i := &Interactor{collectionRepo: cr, groupRepo: gr,
		auth: auth.NewVoidAuth()}
	for _, opt := range opts {
		opt(i)
	}
	return *i

}
