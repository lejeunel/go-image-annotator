package delete

import (
	"fmt"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"log/slog"
)

type Interactor struct {
	repo   Repo
	logger *slog.Logger
}

func (i *Interactor) Execute(r Request, out OutputPort) {
	if err := i.ensureExists(r.Name); err != nil {
		i.handleError(err, out)
		return
	}
	if err := i.ensureDeletable(r.Name); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.repo.Delete(r.Name); err != nil {
		i.handleError(err, out)
		return
	}
	out.Success()
}
func (i *Interactor) ensureDeletable(name string) error {
	baseErr := fmt.Errorf("ensuring collection with name %v is empty", name)
	isPopulated, err := i.repo.IsPopulated(name)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, e.ErrInternal)
	}
	if *isPopulated {
		return fmt.Errorf("%w: %w", baseErr, e.ErrDependency)
	}
	return nil

}

func (i *Interactor) ensureExists(name string) error {
	baseErr := fmt.Errorf("checking whether collection with name %v exists", name)
	exists, err := i.repo.Exists(name)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, e.ErrInternal)
	}
	if !exists {
		return fmt.Errorf("%w: %w", baseErr, e.ErrNotFound)
	}
	return nil
}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "deleting collection"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

func NewInteractor(r Repo) *Interactor {
	return &Interactor{repo: r, logger: logging.NewNoOpLogger()}
}
