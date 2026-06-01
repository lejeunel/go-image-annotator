package add_bbox

import (
	"fmt"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"log/slog"
)

type Interface interface {
	Execute(r Request, out OutputPort)
}

type Interactor struct {
	imageStore st.Interface
	repo       Repo
	logger     *slog.Logger
}

func NewInteractor(imageStore st.Interface, repo Repo) *Interactor {
	return &Interactor{repo: repo, imageStore: imageStore, logger: logging.NewNoOpLogger()}
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

	box := a.NewBoundingBox(a.NewAnnotationId(), r.Xc, r.Yc, r.Width, r.Height, *label)
	if err := i.validateBox(image, box); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.addBox(image, box); err != nil {
		i.handleError(err, out)
		return
	}

	out.SuccessAddBox(box)

}
func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "adding bounding box"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}
func (i *Interactor) addBox(image *im.Image, box a.BoundingBox) error {
	if err := i.repo.AddBoundingBox(image.Id, image.Collection.Id, box); err != nil {
		return err
	}
	return nil
}

func (i *Interactor) validateBox(image *im.Image, box a.BoundingBox) error {
	if err := image.AddBoundingBox(box); err != nil {
		return err
	}
	return nil
}

func (i *Interactor) findLabel(name string) (*lbl.Label, error) {
	label, err := i.repo.FindLabel(name)
	if err != nil {
		return nil, err
	}
	return label, nil
}

func (i *Interactor) findImage(imageId im.ImageId, collectionName string) (*im.Image, error) {
	image, err := i.imageStore.Find(im.BaseImage{ImageId: imageId, Collection: collectionName})
	if err != nil {
		return nil, err
	}
	return image, nil
}
