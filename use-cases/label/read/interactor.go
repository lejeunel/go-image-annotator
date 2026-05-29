package read

import (
	"fmt"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"log/slog"
)

type Interactor struct {
	repo   Repo
	logger *slog.Logger
}

func (i *Interactor) Execute(r Request, out OutputPort) {
	errCtx := "fetching label"
	found, err := i.repo.FindLabelByName(r.Name)

	if err != nil {
		err = fmt.Errorf("%v: %w", errCtx, err)
		i.logger.Error(errCtx, "error", err)
		out.Error(err)
		return
	}

	out.Success(Response{Name: found.Name, Description: found.Description})

}

func NewInteractor(r Repo) *Interactor {
	return &Interactor{repo: r, logger: logging.NewNoOpLogger()}
}
