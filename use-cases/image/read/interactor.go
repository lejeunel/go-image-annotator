package read

import (
	"fmt"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	imstore "github.com/lejeunel/go-image-annotator/modules/image-store"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"log/slog"
)

type Interface interface {
	Execute(Request, OutputPort)
}

type Interactor struct {
	store  imstore.Interface
	logger *slog.Logger
}

func New(store imstore.Interface) Interactor {
	return Interactor{store: store, logger: logging.NewNoOpLogger()}
}

func (i Interactor) Execute(r Request, out OutputPort) {
	errCtx := "reading image meta-data"
	image, err := i.store.Find(im.BaseImage{ImageId: r.ImageId, Collection: r.Collection})
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessReadImage(*image)
}
