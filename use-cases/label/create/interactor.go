package create

import (
	"fmt"

	lbl "github.com/lejeunel/go-image-annotator-v2/entities/label"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
	"github.com/lejeunel/go-image-annotator-v2/shared/logging"
	v "github.com/lejeunel/go-image-annotator-v2/shared/validation"
	"log/slog"
)

type Interactor struct {
	repo      Repo
	validator v.Validator
	logger    *slog.Logger
}

func (i *Interactor) Execute(r Request, out OutputPort) {
	if err := i.validator.Validate(r.Name); err != nil {
		i.handleError(err, out)
		return
	}
	if err := i.checkDuplicate(r.Name); err != nil {
		i.handleError(err, out)
		return
	}

	label := lbl.NewLabel(lbl.NewLabelId(), r.Name, lbl.WithDescription(r.Description))
	if err := i.repo.Create(*label); err != nil {
		i.handleError(err, out)
		return
	}
	out.Success(Response{Name: r.Name, Description: r.Description})
}
func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "creating label"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

func (i *Interactor) checkDuplicate(name string) error {
	errBaseMsg := "checking for duplicate label with name %v: %w"
	alreadyExists, err := i.repo.Exists(name)
	if err != nil {
		return fmt.Errorf(errBaseMsg, name, e.ErrInternal)
	}
	if alreadyExists {
		return fmt.Errorf(errBaseMsg, name, e.ErrDuplicate)
	}
	return nil
}

type Option func(*Interactor)

func WithNameValidator(v v.Validator) Option {
	return func(i *Interactor) {
		i.validator = v
	}
}

func NewInteractor(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r, validator: v.NewNameValidator(),
		logger: logging.NewNoOpLogger()}

	for _, opt := range opts {
		opt(i)
	}
	return i
}
