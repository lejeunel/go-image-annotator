package remove

import (
	"fmt"

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
	errCtx := "deleting bounding box"
	if err := i.repo.RemoveAnnotation(r.Id); err != nil {
		err = fmt.Errorf("%v: %w", errCtx, err)
		i.logger.Error(errCtx, "error", err)
		out.Error(err)
		return
	}

	out.SuccessDeleteAnnotation(Response{Id: r.Id})

}
