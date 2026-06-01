package update

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

func NewInteractor(r Repo) Interactor {
	return Interactor{repo: r, logger: logging.NewNoOpLogger()}
}

func (i *Interactor) Execute(r Request, out OutputPort) {

	if err := i.ensureNameExists(r.Name); err != nil {
		i.handleError(err, out)
		return
	}

	if r.NewName != r.Name {
		if err := i.ensureNameDoesNotExist(r.NewName); err != nil {
			i.handleError(err, out)
			return
		}

	}

	if err := i.repo.Update(Model{Name: r.Name, NewName: r.NewName, NewDescription: r.NewDescription}); err != nil {
		i.handleError(err, out)
		return
	}

	out.Success(Response{Name: r.NewName, Description: r.NewDescription})
}

func (i *Interactor) ensureNameExists(name string) error {
	baseErr := fmt.Errorf("ensuring that collection with name %v exists", name)
	exists, err := i.repo.Exists(name)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, e.ErrInternal)
	}
	if !exists {
		return fmt.Errorf("%w: %w", baseErr, e.ErrNotFound)
	}
	return nil
}

func (i *Interactor) ensureNameDoesNotExist(name string) error {
	baseErr := fmt.Errorf("ensuring that a collection with name %v does not already exist", name)
	exists, err := i.repo.Exists(name)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, e.ErrInternal)
	}
	if exists {
		return fmt.Errorf("%w: %w", baseErr, e.ErrDuplicate)
	}
	return nil
}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "updating collection"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}
