package import_image

import (
	"fmt"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"log/slog"
)

type Interactor struct {
	repo   Repo
	logger *slog.Logger
}

func NewInteractor(repo Repo) *Interactor {
	return &Interactor{repo: repo, logger: logging.NewNoOpLogger()}
}
func (i *Interactor) Execute(r Request, out OutputPort) {

	if err := i.ensureSourceImageExists(r.ImageId); err != nil {
		i.handleError(err, out)
		return
	}
	collection, err := i.findDestinationCollection(r.Collection)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.ensureImageDoesNotAlreadyExistInCollection(r.ImageId, collection.Id); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.repo.AddToCollection(r.ImageId, collection.Id); err != nil {
		i.handleError(err, out)
		return
	}

	out.Success(Response{})

}
func (i *Interactor) ensureImageDoesNotAlreadyExistInCollection(imageId im.ImageId, collectionId clc.CollectionId) error {

	baseErr := fmt.Errorf("ensuring that source image does not already exist in destination collection")
	alreadyExists, err := i.repo.ImageExistsInCollection(imageId, collectionId)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	if alreadyExists {
		return fmt.Errorf("%w: %w", baseErr, e.ErrDependency)
	}
	return nil
}

func (i *Interactor) ensureSourceImageExists(id im.ImageId) error {
	baseErr := fmt.Errorf("ensuring that source image exists")
	exists, err := i.repo.ImageExists(id)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	if !exists {
		return fmt.Errorf("%w: %w", baseErr, e.ErrNotFound)
	}
	return nil

}

func (i *Interactor) findDestinationCollection(name string) (*clc.Collection, error) {

	baseErr := fmt.Errorf("fetching destination collection")
	collection, err := i.repo.FindCollectionByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return collection, nil

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "deleting image"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}
