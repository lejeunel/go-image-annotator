package read

import (
	"context"
	"fmt"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Interactor struct {
	CollectionRepo
	ImageRepo
	MetaDataRepo
}

func New(c CollectionRepo, ir ImageRepo, m MetaDataRepo) Interactor {
	return Interactor{c, ir, m}
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := fmt.Errorf("reading metadata for image %v, collection %v and key %v",
		r.ImageId, r.Collection, r.Key)

	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: parsing image id %v: %w", errCtx, imageId, err))
		return
	}

	collectionExists, err := i.CollectionRepo.Exists(r.Collection)
	if err != nil {
		out.Error(
			fmt.Errorf(
				"%v: checking existence of collection: %v: %w",
				errCtx,
				err,
				e.ErrInternal,
			),
		)
		return
	}
	if !collectionExists {
		out.Error(
			fmt.Errorf(
				"%v: checking existence of collection: %v: %w",
				errCtx,
				err,
				e.ErrValidation,
			),
		)
		return
	}

	imageInCollection, err := i.ImageRepo.ImageExistsInCollection(imageId, r.Collection)
	if err != nil {
		out.Error(
			fmt.Errorf(
				"%v: checking whether image is in collection: %v: %w",
				errCtx,
				err,
				e.ErrInternal,
			),
		)
		return
	}
	if !imageInCollection {
		out.Error(
			fmt.Errorf(
				"%v: checking whether image is in collection: %v: %w",
				errCtx,
				err,
				e.ErrValidation,
			),
		)
		return
	}

	exists, err := i.MetaDataRepo.KeyExists(r.Collection, imageId, r.Key)
	if err != nil {
		out.Error(
			fmt.Errorf(
				"%v: checking existence of key: %v: %w",
				errCtx,
				err,
				e.ErrInternal,
			),
		)
		return
	}
	if !exists {
		out.Error(
			fmt.Errorf("%v: checking existence of key: %w", errCtx, e.ErrValidation),
		)
		return
	}

	value, err := i.MetaDataRepo.GetValue(r.Collection, imageId, r.Key)
	if err != nil {
		out.Error(
			fmt.Errorf(
				"%v: fetching value: %v: %w",
				errCtx,
				err,
				e.ErrInternal,
			),
		)
		return
	}

	out.SuccessReadMetadata(Response{
		ImageId:    r.ImageId,
		Collection: r.Collection,
		Data:       m.MetaData{Key: r.Key, Value: *value},
	})
}
