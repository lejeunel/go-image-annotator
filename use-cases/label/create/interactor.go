package create

import (
	"context"
	"fmt"

	"log/slog"

	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	"github.com/lejeunel/go-image-annotator/shared/auth"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	v "github.com/lejeunel/go-image-annotator/shared/validation"
)

type Interactor struct {
	repo      Repo
	validator v.Validator
	logger    *slog.Logger
	auth      Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "creating label"
	if err := i.auth.CreateLabel(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return

	}
	if err := i.validator.Validate(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if err := i.checkDuplicate(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	label := lbl.NewLabel(lbl.NewLabelId(), r.Name, lbl.WithDescription(r.Description))
	if err := i.repo.Create(label); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.Success(Response{Name: r.Name, Description: r.Description})
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
func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r, validator: v.NewNameValidator(),
		auth: auth.PassThroughAuth{},
	}

	for _, opt := range opts {
		opt(i)
	}
	return i
}
