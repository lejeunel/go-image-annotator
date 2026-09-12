package create

import (
	"context"
	"fmt"

	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	v "github.com/lejeunel/go-image-annotator/modules/string-validator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	Repo
	v.Validator
	Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "creating profile"
	if err := i.Auth.CreateProfile(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	errBaseMsg := fmt.Sprintf("checking for duplicate profile with name %v", r.Name)
	alreadyExists, err := i.Repo.Exists(r.Name)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if *alreadyExists {
		out.Error(fmt.Errorf("%v: %w", errBaseMsg, e.ErrDuplicate))
		return
	}
	if err := i.Validator.Validate(r.Name); err != nil {
		out.Error(fmt.Errorf("validating profile name %v: %w", r.Name, err))
		return
	}

	profile := pr.NewProfile(
		pr.NewProfileId(),
		r.Name,
		pr.WithDescription(r.Description),
		pr.WithLabels(r.Labels))

	if err := i.Repo.Create(profile); err != nil {
		out.Error(fmt.Errorf("creating profile with name %v: %w", r.Name, err))
		return
	}
	out.SuccessCreateProfile(profile)

}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}
func WithNameValidator(v v.Validator) Option {
	return func(i *Interactor) {
		i.Validator = v
	}
}

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{
		Repo: r, Validator: v.NewNameValidator(),
		Auth: auth.NewVoidAuth(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
