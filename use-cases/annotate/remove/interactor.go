package remove

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
	id, err := a.NewAnnotationIdFromString(r.Id)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.repo.RemoveAnnotation(*id); err != nil {
		i.handleError(err, out)
		return
	}

	out.SuccessDeleteAnnotation(Response{Id: *id})

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "removing annotation"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}
