package find

import (
	"fmt"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	imstore "github.com/lejeunel/go-image-annotator/modules/image-store"
)

type Interface interface {
	Execute(Request, OutputPort)
}

type Interactor struct {
	store imstore.Interface
}

func New(store imstore.Interface) Interactor {
	return Interactor{store: store}
}

func (i Interactor) Execute(r Request, out OutputPort) {
	errCtx := "reading image meta-data"
	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	image, err := i.store.Find(im.BaseImage{ImageId: imageId, Collection: r.Collection})
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessReadImage(*image)
}
