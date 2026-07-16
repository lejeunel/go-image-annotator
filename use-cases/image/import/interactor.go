package import_image

import (
	"context"
	"fmt"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	ImageRepo
	CollectionRepo
	auth Auth
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(imr ImageRepo, c CollectionRepo, opts ...Option) *Interactor {
	i := &Interactor{ImageRepo: imr, CollectionRepo: c,
		auth: auth.NewVoidAuth()}
	for _, opt := range opts {
		opt(i)
	}
	return i

}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "importing image"
	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.ensureSourceImageExists(imageId); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	dstCollection, err := i.findCollection(r.DestinationCollection)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if dstCollection.Group != nil {
		if err := i.auth.ImportImage(ctx, dstCollection.Group.Name); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	if err := i.ensureImageDoesNotAlreadyExistInCollection(imageId, dstCollection.Id); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.ImageRepo.AddToCollection(imageId, dstCollection.Id); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	// TODO import annotations here

	out.Success(Response{})

}
func (i Interactor) ensureImageDoesNotAlreadyExistInCollection(imageId im.ImageId, collectionId clc.CollectionId) error {

	errCtx := fmt.Errorf("ensuring that source image does not already exist in destination collection")
	alreadyExists, err := i.ImageRepo.ImageExistsInCollection(imageId, collectionId)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, err)
	}
	if alreadyExists {
		return fmt.Errorf("%w: %w", errCtx, e.ErrDependency)
	}
	return nil
}
func (i Interactor) ensureSourceImageExists(id im.ImageId) error {
	errCtx := fmt.Errorf("ensuring that source image exists")
	exists, err := i.ImageRepo.ImageExists(id)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, err)
	}
	if !exists {
		return fmt.Errorf("%w: %w", errCtx, e.ErrNotFound)
	}
	return nil

}
func (i Interactor) findCollection(name string) (*clc.Collection, error) {

	errCtx := fmt.Errorf("fetching collection %v", name)
	collection, err := i.CollectionRepo.FindCollectionByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errCtx, err)
	}
	return collection, nil

}
