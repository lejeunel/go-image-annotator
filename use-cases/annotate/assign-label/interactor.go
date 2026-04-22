package assign_label

import (
	"fmt"

	st "github.com/lejeunel/go-image-annotator-v2/app/image-store"
	clc "github.com/lejeunel/go-image-annotator-v2/entities/collection"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	lbl "github.com/lejeunel/go-image-annotator-v2/entities/label"
	"github.com/lejeunel/go-image-annotator-v2/shared/logging"
	"log/slog"
)

type Interactor struct {
	repo   Repo
	store  st.Interface
	logger *slog.Logger
}

func (i *Interactor) Execute(r Request, out OutputPort) {

	image, err := i.findImage(r.ImageId, r.Collection)
	if err != nil {
		i.handleError(err, out)
		return
	}
	label, err := i.findLabel(r.Label)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.addLabel(image.Id, image.Collection.Id, label.Id); err != nil {
		i.handleError(err, out)
		return
	}

	out.Success(Response{ImageId: r.ImageId, Collection: r.Collection, Label: r.Label})
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

func (i *Interactor) addLabel(imageId im.ImageId, collectionId clc.CollectionId, labelId lbl.LabelId) error {
	if err := i.repo.AddImageLabel(imageId, collectionId, labelId); err != nil {
		return err
	}
	return nil

}

func NewInteractor(repo Repo, store st.Interface) *Interactor {
	return &Interactor{repo: repo, store: store, logger: logging.NewNoOpLogger()}
}
