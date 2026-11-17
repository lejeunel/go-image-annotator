package copy

import (
	"context"
	"errors"
	"fmt"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type ImageStore interface {
	Find(im.BaseImage) (*im.Image, error)
	DeepCopy(src clc.CollectionName, id im.ImageId, dst clc.CollectionName) error
	ShallowCopy(id im.ImageId, dst clc.CollectionName) error
}

type CollectionRepo interface {
	Find(string) (*clc.Collection, error)
}

type Interactor struct {
	ImageStore
	CollectionRepo
	Auth
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(store ImageStore, collections CollectionRepo, opts ...Option) Interactor {
	i := &Interactor{
		ImageStore:     store,
		CollectionRepo: collections,
		Auth:           auth.NewVoidAuth(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := fmt.Sprintf("copying image %v from collection %v to collection %v",
		r.ImageId, r.Source, r.Destination)

	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if _, err := i.ImageStore.Find(
		im.BaseImage{ImageId: imageId, Collection: r.Source},
	); err != nil {
		out.Error(fmt.Errorf("%v: finding source image: %w", errCtx, err))
		return
	}

	dst, err := i.CollectionRepo.Find(r.Destination)
	if err != nil {
		out.Error(fmt.Errorf("%v: finding destination collection: %w", errCtx, err))
		return
	}

	if err := i.Auth.CopyImage(ctx, dst.Group); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.checkNotInDestination(imageId, dst.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if r.Deep {
		err = i.ImageStore.DeepCopy(r.Source, imageId, dst.Name)
	} else {
		err = i.ImageStore.ShallowCopy(imageId, dst.Name)
	}
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	copied, err := i.ImageStore.Find(im.BaseImage{ImageId: imageId, Collection: dst.Name})
	if err != nil {
		out.Error(fmt.Errorf("%v: reading back copied image: %w", errCtx, err))
		return
	}
	out.SuccessCopyImage(Response{Image: *copied})
}

func (i *Interactor) checkNotInDestination(id im.ImageId, dst clc.CollectionName) error {
	_, err := i.ImageStore.Find(im.BaseImage{ImageId: id, Collection: dst})
	switch {
	case err == nil:
		return fmt.Errorf("image already in collection %v: %w", dst, e.ErrDuplicate)
	case errors.Is(err, e.ErrNotFound):
		return nil
	default:
		return fmt.Errorf("checking for image in destination collection: %w", err)
	}
}
