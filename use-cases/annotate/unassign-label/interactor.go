package unassign_label

import (
	"fmt"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"log/slog"
)

type Interactor struct {
	repo   Repo
	output OutputPort
	store  st.Interface
	logger *slog.Logger
}

func (i *Interactor) Execute(r Request, out OutputPort) {
	image, err := i.findImage(r.ImageId, r.Collection)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.removeLabel(*image, r.Label); err != nil {
		i.handleError(err, out)
		return
	}

	out.Success(Response{})

}

func (i *Interactor) findImage(imageId im.ImageId, collection string) (*im.Image, error) {
	baseErr := fmt.Errorf("finding image %v in collection %v", imageId, collection)
	image, err := i.store.Find(im.BaseImage{ImageId: imageId, Collection: collection})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return image, nil
}

func (i *Interactor) removeLabel(image im.Image, label string) error {
	baseErr := fmt.Errorf("removing image label")
	removed := 0
	for _, imageLabel := range image.Labels {
		if imageLabel.Label.Name == label {
			err := i.repo.RemoveImageLabel(image.Id, image.Collection.Id, imageLabel.Label.Id)
			if err != nil {
				return fmt.Errorf("%w: %w", baseErr, err)
			}
			removed += 1
		}
	}

	if removed == 0 {
		return fmt.Errorf("%w: %w", baseErr, e.ErrNotFound)
	}
	return nil

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "deleting image label"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

func NewInteractor(repo Repo, store st.Interface) *Interactor {
	return &Interactor{repo: repo, store: store, logger: logging.NewNoOpLogger()}
}
