package delete

import (
	"context"
	"fmt"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	st "github.com/lejeunel/go-image-annotator/modules/image-store"
)

type Interactor struct {
	store st.Interface
	ImageRepo
	AnnotationRepo
	auth Auth
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(store st.Interface, r ImageRepo, a AnnotationRepo, opts ...Option) *Interactor {
	i := &Interactor{store: store, ImageRepo: r,
		AnnotationRepo: a,
		auth:           auth.NewVoidAuth()}
	for _, opt := range opts {
		opt(i)
	}
	return i
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
		if err := i.auth.DeleteImage(ctx, image.Collection.Group.Name); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	if err := i.deleteLabels(*image); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.deleteBoundingBoxes(*image); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.ImageRepo.RemoveImageFromCollection(image.Id, image.Collection.Id); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return

	}

	out.SuccessDeleteImage(Response{})
}

func (i *Interactor) deleteBoundingBoxes(image im.Image) error {
	baseErr := fmt.Errorf("deleting bounding box annotations")
	for _, box := range image.BoundingBoxes {
		if err := i.AnnotationRepo.RemoveAnnotation(box.Id); err != nil {
			return fmt.Errorf("%w: %w", baseErr, err)
		}
	}
	return nil

}

func (i *Interactor) deleteLabels(image im.Image) error {
	baseErr := fmt.Errorf("deleting image labels")
	for _, label := range image.Labels {
		if err := i.AnnotationRepo.RemoveAnnotation(label.Id); err != nil {
			return fmt.Errorf("%w: %w", baseErr, err)
		}
	}
	return nil

}

func (i *Interactor) findImage(imageId im.ImageId, collection string) (*im.Image, error) {
	baseErr := fmt.Errorf("fetching associated resources")
	image, err := i.store.Find(im.BaseImage{ImageId: imageId, Collection: collection})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return image, nil

}
