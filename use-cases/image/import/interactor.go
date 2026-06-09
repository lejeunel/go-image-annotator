package import_image

import (
	"context"
	"fmt"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"github.com/lejeunel/go-image-annotator/shared/auth"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"log/slog"
)

type Interactor struct {
	repo   Repo
	logger *slog.Logger
	auth   Auth
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func NewInteractor(repo Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: repo, logger: logging.NewNoOpLogger(),
		auth: auth.PassThroughAuth{}}
	for _, opt := range opts {
		opt(i)
	}
	return i

}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.ensureSourceImageExists(imageId); err != nil {
		i.handleError(err, out)
		return
	}
	dstCollection, err := i.findCollection(r.DestinationCollection)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if dstCollection.Group != nil {
		if err := i.auth.ImportImage(ctx, dstCollection.Group.Name); err != nil {
			i.handleError(err, out)
			return
		}
	}

	if err := i.ensureImageDoesNotAlreadyExistInCollection(imageId, dstCollection.Id); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.repo.AddToCollection(imageId, dstCollection.Id); err != nil {
		i.handleError(err, out)
		return
	}

	// TODO import annotations here

	out.Success(Response{})

}
func (i *Interactor) ensureImageDoesNotAlreadyExistInCollection(imageId im.ImageId, collectionId clc.CollectionId) error {

	errCtx := fmt.Errorf("ensuring that source image does not already exist in destination collection")
	alreadyExists, err := i.repo.ImageExistsInCollection(imageId, collectionId)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, err)
	}
	if alreadyExists {
		return fmt.Errorf("%w: %w", errCtx, e.ErrDependency)
	}
	return nil
}

func (i *Interactor) ensureSourceImageExists(id im.ImageId) error {
	errCtx := fmt.Errorf("ensuring that source image exists")
	exists, err := i.repo.ImageExists(id)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, err)
	}
	if !exists {
		return fmt.Errorf("%w: %w", errCtx, e.ErrNotFound)
	}
	return nil

}

func (i *Interactor) findCollection(name string) (*clc.Collection, error) {

	errCtx := fmt.Errorf("fetching collection %v", name)
	collection, err := i.repo.FindCollectionByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errCtx, err)
	}
	return collection, nil

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "deleting image"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}
