package update_label

import (
	"fmt"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
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
	label, err := i.repo.FindLabelByName(r.Label)
	if err != nil {
		i.handleError(err, out)
		return
	}

	id, err := a.NewAnnotationIdFromString(r.AnnotationId)
	if err != nil {
		i.handleError(err, out)
	}

	err = i.repo.UpdateLabelOfAnnotation(*id, label.Id)
	if err != nil {
		i.handleError(err, out)
		return
	}

	out.SuccessUpdateLabel(Response{})

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "updating bounding box properties"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}
