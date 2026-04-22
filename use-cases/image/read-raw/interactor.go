package read_raw

import (
	"fmt"

	ast "github.com/lejeunel/go-image-annotator-v2/app/file-store"
	"github.com/lejeunel/go-image-annotator-v2/shared/logging"
	"log/slog"
)

type Interactor struct {
	repo   ast.ReadInterface
	logger *slog.Logger
}

func (i *Interactor) Execute(r Request, out OutputPort) {
	reader, err := i.repo.Get(r.ImageId)
	if err != nil {
		i.handleError(err, out)
		return
	}

	out.Success(Response{Reader: reader})

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "reading raw image data"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

func NewInteractor(repo ast.ReadInterface) *Interactor {
	return &Interactor{repo: repo, logger: logging.NewNoOpLogger()}
}
