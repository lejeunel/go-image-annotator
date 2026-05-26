package read

import (
	"fmt"
	imstore "github.com/lejeunel/go-image-annotator-v2/app/image-store"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	"github.com/lejeunel/go-image-annotator-v2/shared/logging"
	"log/slog"
)

type Interface interface {
	Execute(Request, OutputPort)
}

type Interactor struct {
	store  imstore.Interface
	logger *slog.Logger
}

func NewInteractor(store imstore.Interface) *Interactor {
	return &Interactor{store: store, logger: logging.NewNoOpLogger()}
}

func (i *Interactor) Execute(r Request, out OutputPort) {
	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		i.handleError(err, out)
		return
	}

	image, err := i.store.Find(im.BaseImage{ImageId: imageId, Collection: r.Collection})
	if err != nil {
		i.handleError(err, out)
		return
	}

	out.SuccessReadImage(*image)
}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "reading image meta-data"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}
