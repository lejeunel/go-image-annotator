package create

import (
	"context"
	"fmt"

	"github.com/jonboulle/clockwork"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	v "github.com/lejeunel/go-image-annotator/shared/validation"
)

type Interactor struct {
	collectionRepo CollectionRepo
	groupRepo      GroupRepo
	validator      v.Validator
	clock          clockwork.Clock
	auth           Auth
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "creating collection"
	if r.Group != nil {
		if err := i.auth.CreateCollection(ctx, *r.Group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}
	if err := i.validate(r.Name); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.create(r); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.Success(Response{Name: r.Name, Description: r.Description})
}
func (i Interactor) create(r Request) error {
	collection := clc.NewCollection(clc.NewCollectionId(), r.Name,
		clc.WithDescription(r.Description),
		clc.WithCreatedAt(i.clock.Now()))
	if r.Group != nil {
		group, err := i.groupRepo.Find(*r.Group)
		if err != nil {
			return err
		}
		collection.Group = group
	}
	if err := i.collectionRepo.Create(collection); err != nil {
		return err
	}
	return nil

}
func (i Interactor) validate(name string) error {
	if err := i.validator.Validate(name); err != nil {
		return fmt.Errorf("checking collection name %v: %w", name, err)
	}
	if err := i.isDuplicate(name); err != nil {
		return err
	}
	return nil

}
func (i Interactor) isDuplicate(name string) error {
	errBaseMsg := fmt.Sprintf("checking for duplicate collection with name %v", name)
	alreadyExists, err := i.collectionRepo.Exists(name)
	if err != nil {
		return fmt.Errorf("%v: %w", errBaseMsg, e.ErrInternal)
	}
	if alreadyExists {
		return fmt.Errorf("%v: %w", errBaseMsg, e.ErrDuplicate)
	}
	return nil
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

func New(rc CollectionRepo, rg GroupRepo, opts ...Option) Interactor {
	i := &Interactor{collectionRepo: rc,
		groupRepo: rg,
		validator: v.NewNameValidator(),
		clock:     clockwork.NewRealClock(),
		auth:      auth.NewVoidAuth()}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
