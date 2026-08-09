package delete

import (
	"context"
	"fmt"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
)

type ImageStore interface {
	Find(im.BaseImage) (*im.Image, error)
	Delete(im.ImageId, clc.CollectionName) error
}

type Interactor struct {
	ImageStore
	auth Auth
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(store ImageStore, opts ...Option) Interactor {
	i := &Interactor{
		ImageStore: store,
		auth:       auth.NewVoidAuth(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "deleting image"

	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
	}

	image, err := i.findImage(imageId, r.Collection)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if image.Collection.Group != nil {
		if err := i.auth.DeleteImage(ctx, *image.Collection.Group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	if err := i.ImageStore.Delete(imageId, r.Collection); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessDeleteImage(Response{ImageId: image.Id.String(), Collection: image.Collection.Name})
}

func (i *Interactor) findImage(imageId im.ImageId, collection string) (*im.Image, error) {
	baseErr := fmt.Errorf("fetching associated resources")
	image, err := i.ImageStore.Find(im.BaseImage{ImageId: imageId, Collection: collection})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return image, nil
}
