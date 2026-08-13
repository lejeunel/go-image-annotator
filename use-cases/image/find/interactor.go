package find

import (
	"fmt"

	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Interface interface {
	Execute(Request, OutputPort)
}
type ImageStore interface {
	Find(base im.BaseImage) (*im.Image, error)
}

type Interactor struct {
	ImageStore
}

func New(store ImageStore) Interactor {
	return Interactor{ImageStore: store}
}

func (i Interactor) Execute(r Request, out OutputPort) {
	errCtx := "reading image meta-data"
	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	image, err := i.ImageStore.Find(im.BaseImage{ImageId: imageId, Collection: r.Collection})
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessReadImage(*image)
}
