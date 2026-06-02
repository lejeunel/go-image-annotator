package delete

import (
	"fmt"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"log/slog"
)

type Interactor struct {
	store  st.Interface
	repo   Repo
	logger *slog.Logger
}

func NewInteractor(store st.Interface, repo Repo) *Interactor {
	return &Interactor{store: store, repo: repo, logger: logging.NewNoOpLogger()}
}

func (i *Interactor) Execute(r Request, out OutputPort) {
	image, err := i.findImage(r.ImageId, r.Collection)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.deleteLabels(*image); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.deleteBoundingBoxes(*image); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.repo.RemoveImageFromCollection(image.Id, image.Collection.Id); err != nil {
		i.handleError(err, out)
		return

	}

	out.Success(Response{})
}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "deleting image"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

func (i *Interactor) deleteBoundingBoxes(image im.Image) error {
	baseErr := fmt.Errorf("deleting bounding box annotations")
	for _, box := range image.BoundingBoxes {
		if err := i.repo.RemoveAnnotation(image.Id, image.Collection.Id, box.Id); err != nil {
			return fmt.Errorf("%w: %w", baseErr, err)
		}
	}
	return nil

}

func (i *Interactor) deleteLabels(image im.Image) error {
	baseErr := fmt.Errorf("deleting image labels")
	for _, label := range image.Labels {
		if err := i.repo.RemoveAnnotation(image.Id, image.Collection.Id, label.Id); err != nil {
			return fmt.Errorf("%w: %w", baseErr, err)
		}
	}
	return nil

}

func (i *Interactor) findImage(imageId string, collection string) (*im.Image, error) {
	baseErr := fmt.Errorf("fetching associated resources")
	image, err := i.store.Find(im.BaseImage{ImageId: imageId, Collection: collection})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return image, nil

}
