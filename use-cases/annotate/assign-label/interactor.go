package assign_label

import (
	"fmt"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"log/slog"
)

type Interface interface {
	Execute(r Request, out OutputPort)
}

type Interactor struct {
	repo   Repo
	store  st.Interface
	logger *slog.Logger
}

func (i *Interactor) Execute(r Request, out OutputPort) {

	id, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		i.handleError(err, out)
		return
	}

	image, err := i.findImage(id, r.Collection)
	if err != nil {
		i.handleError(err, out)
		return
	}
	label, err := i.findLabel(r.Label)
	if err != nil {
		i.handleError(err, out)
		return
	}

	imageLabel, err := i.addLabel(image.Id, image.Collection.Id, *label)
	if err != nil {
		i.handleError(err, out)
		return
	}

	out.SuccessAddLabel(Response{
		ImageId:      r.ImageId,
		Collection:   r.Collection,
		Label:        r.Label,
		AnnotationId: imageLabel.Id.String()})
}
func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "assigning label to image"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)

}
func (i *Interactor) findLabel(name string) (*lbl.Label, error) {
	label, err := i.repo.FindLabel(name)
	if err != nil {
		return nil, err
	}
	return label, nil

}
func (i *Interactor) findImage(imageId im.ImageId, collection string) (*im.Image, error) {
	image, err := i.store.Find(im.BaseImage{ImageId: imageId, Collection: collection})
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (i *Interactor) addLabel(imageId im.ImageId, collectionId clc.CollectionId, label lbl.Label) (*an.ImageLabel, error) {
	imageLabel := an.NewImageLabel(label)
	if err := i.repo.AddImageLabel(imageId, collectionId, imageLabel); err != nil {
		return nil, err
	}
	return &imageLabel, nil

}

func NewInteractor(repo Repo, store st.Interface) *Interactor {
	return &Interactor{repo: repo, store: store, logger: logging.NewNoOpLogger()}
}
