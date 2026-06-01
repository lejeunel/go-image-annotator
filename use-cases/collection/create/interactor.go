package create

import (
	"fmt"
	"log/slog"

	"github.com/jonboulle/clockwork"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	auth "github.com/lejeunel/go-image-annotator/shared/auth"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	v "github.com/lejeunel/go-image-annotator/shared/validation"
)

type Interactor struct {
	repo      Repo
	validator v.Validator
	logger    *slog.Logger
	clock     clockwork.Clock
	auth      Auth
}

func (i *Interactor) Execute(p auth.PrincipalProvider, r Request, out OutputPort) {
	if err := i.auth.CreateCollection(p, r.Group); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.validate(r.Name); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.create(r); err != nil {
		i.handleError(err, out)
		return
	}

	out.Success(Response{Name: r.Name, Description: r.Description})
}

func (i *Interactor) create(r Request) error {
	collection := clc.NewCollection(clc.NewCollectionId(), r.Name,
		clc.WithDescription(r.Description),
		clc.WithCreatedAt(i.clock.Now()))
	if err := i.repo.Create(collection); err != nil {
		return err
	}
	return nil

}

func (i *Interactor) validate(name string) error {
	if err := i.validator.Validate(name); err != nil {
		return fmt.Errorf("checking collection name %v: %w", name, err)
	}
	if err := i.isDuplicate(name); err != nil {
		return err
	}
	return nil

}

func (i *Interactor) isDuplicate(name string) error {
	errBaseMsg := fmt.Sprintf("checking for duplicate collection with name %v", name)
	alreadyExists, err := i.repo.Exists(name)
	if err != nil {
		return fmt.Errorf("%v: %w", errBaseMsg, e.ErrInternal)
	}
	if alreadyExists {
		return fmt.Errorf("%v: %w", errBaseMsg, e.ErrDuplicate)
	}
	return nil
}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "creating collection"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

type Option func(*Interactor)

func WithNameValidator(v v.Validator) Option {
	return func(i *Interactor) {
		i.validator = v
	}
}

func WithClock(c clockwork.Clock) Option {
	return func(i *Interactor) {
		i.clock = c
	}
}

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func NewInteractor(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r, validator: v.NewNameValidator(),
		logger: logging.NewNoOpLogger(),
		clock:  clockwork.NewRealClock(),
		auth:   auth.PassThroughAuth{}}

	for _, opt := range opts {
		opt(i)
	}
	return i
}
