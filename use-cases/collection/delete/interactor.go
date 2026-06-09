package delete

import (
	"context"
	"fmt"
	auth "github.com/lejeunel/go-image-annotator/shared/auth"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"log/slog"
)

type Interactor struct {
	collectionRepo CollectionRepo
	groupRepo      GroupRepo
	logger         *slog.Logger
	auth           Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	if err := i.authorizeDeletion(ctx, r.Name); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.ensureExists(r.Name); err != nil {
		i.handleError(err, out)
		return
	}
	if err := i.ensureDeletable(r.Name); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.collectionRepo.Delete(r.Name); err != nil {
		i.handleError(err, out)
		return
	}
	out.Success()
}

func (i *Interactor) authorizeDeletion(ctx context.Context, name string) error {
	errCtx := fmt.Errorf("checking group ownership of collection with name %v is empty", name)
	group, err := i.groupRepo.GroupOfCollection(name)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, e.ErrInternal)
	}
	if group != nil {
		if err := i.auth.DeleteCollection(ctx, *group); err != nil {
			return fmt.Errorf("%w: %w", errCtx, e.ErrAuth)
		}
	}
	return nil

}

func (i *Interactor) ensureDeletable(name string) error {
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

func (i *Interactor) ensureExists(name string) error {
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

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "deleting collection"
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

func NewInteractor(cr CollectionRepo, gr GroupRepo, opts ...Option) Interactor {
	i := &Interactor{collectionRepo: cr, groupRepo: gr,
		logger: logging.NewNoOpLogger(),
		auth:   auth.PassThroughAuth{}}
	for _, opt := range opts {
		opt(i)
	}
	return *i

}
