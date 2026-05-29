package modify_bbox

import (
	"fmt"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"log/slog"
)

type Interface interface {
	Execute(Request, OutputPort)
}

type Interactor struct {
	repo   Repo
	logger *slog.Logger
}

func NewInteractor(repo Repo) *Interactor {
	return &Interactor{repo: repo, logger: logging.NewNoOpLogger()}
}
func (i *Interactor) Execute(r Request, out OutputPort) {
	label, err := i.findLabel(r.Label)
	if err != nil {
		i.handleError(err, out)
		return
	}
	u, err := i.validate(r.Xc, r.Yc, r.Width, r.Height, *label)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.update(r.AnnotationId, *u); err != nil {
		i.handleError(err, out)
		return
	}
	out.SuccessUpdateBox(Response{})

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "updating bounding box properties"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

func (i *Interactor) update(id a.AnnotationId, u a.BoundingBoxUpdatables) error {
	if err := i.repo.UpdateBoundingBox(id, u); err != nil {
		return err
	}
	return nil

}

func (i *Interactor) validate(xc float32, yc float32, width float32,
	height float32, label lbl.Label) (*a.BoundingBoxUpdatables, error) {

	if err := a.ValidateBoundingBox(xc, yc, width, height); err != nil {
		return nil, err
	}
	return &a.BoundingBoxUpdatables{LabelId: label.Id, Xc: xc, Yc: yc, Width: width, Height: height}, nil

}

func (i *Interactor) findLabel(name string) (*lbl.Label, error) {

	label, err := i.repo.FindLabelByName(name)
	if err != nil {
		return nil, err
	}
	return label, nil

}
