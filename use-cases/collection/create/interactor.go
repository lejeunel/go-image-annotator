package create

import (
	"context"
	"fmt"

	"github.com/jonboulle/clockwork"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	v "github.com/lejeunel/go-image-annotator/modules/string-validator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	CollectionRepo
	GroupRepo
	ProfileRepo
	v.Validator
	clockwork.Clock
	Auth
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "creating collection"
	if err := i.Auth.CreateCollection(ctx, r.Group); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
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
		clc.WithCreatedAt(i.Clock.Now()))
	if r.Group != nil {
		group, err := i.GroupRepo.Find(*r.Group)
		if err != nil {
			return err
		}
		collection.Group = &group.Name
	}
	if r.Profile != nil {
		profile, err := i.ProfileRepo.Find(*r.Profile)
		if err != nil {
			return err
		}
		collection.Profile = &profile.Name
	}
	if err := i.CollectionRepo.Create(collection); err != nil {
		return err
	}
	return nil
}

func (i Interactor) validate(name string) error {
	if err := i.Validator.Validate(name); err != nil {
		return fmt.Errorf("checking collection name %v: %w", name, err)
	}
	if err := i.isDuplicate(name); err != nil {
		return err
	}
	return nil
}

func (i Interactor) isDuplicate(name string) error {
	errBaseMsg := fmt.Sprintf("checking for duplicate collection with name %v", name)
	alreadyExists, err := i.CollectionRepo.Exists(name)
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
		i.Validator = v
	}
}

func WithClock(c clockwork.Clock) Option {
	return func(i *Interactor) {
		i.Clock = c
	}
}

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(rc CollectionRepo, rg GroupRepo, pr ProfileRepo, opts ...Option) Interactor {
	i := &Interactor{
		CollectionRepo: rc,
		GroupRepo:      rg,
		ProfileRepo:    pr,
		Validator:      v.NewNameValidator(),
		Clock:          clockwork.NewRealClock(),
		Auth:           auth.NewVoidAuth(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
