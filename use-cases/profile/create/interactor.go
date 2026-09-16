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
	ProfileRepo
	LabelRepo
	GroupRepo
	v.Validator
	Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := fmt.Errorf("creating profile with name %v", r.Name)
	if err := i.Auth.CreateProfile(ctx, r.Group); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	errBase := fmt.Errorf("%w: checking for duplicate profile name %v", errCtx, r.Name)
	alreadyExists, err := i.ProfileRepo.Exists(r.Name)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if *alreadyExists {
		out.Error(fmt.Errorf("%w: %w", errBase, e.ErrDuplicate))
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
		pr.WithLabels(r.Labels),
	)

	if r.Group != nil {
		exists, err := i.GroupRepo.Exists(*r.Group)
		if err != nil {
			out.Error(fmt.Errorf("%w: checking existence of group %v: %w", errCtx, *r.Group, err))
			return
		}
		if !*exists {
			out.Error(
				fmt.Errorf(
					"%w: checking existence of group %v: %w",
					errCtx,
					*r.Group,
					e.ErrValidation,
				),
			)
			return
		}
		profile.Group = r.Group
	}

	for _, label := range profile.Labels {
		errBase := fmt.Errorf("%w: checking for existence of label %v", errCtx, label)
		exists, err := i.LabelRepo.Exists(label)
		if err != nil {
			out.Error(fmt.Errorf("%w: %w", errBase, err))
			return
		}
		if !exists {
			out.Error(fmt.Errorf("%w: %w", errBase, e.ErrValidation))
			return
		}
	}

	if err := i.ProfileRepo.Create(profile); err != nil {
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

func New(r ProfileRepo, l LabelRepo, g GroupRepo, opts ...Option) Interactor {
	i := &Interactor{
		ProfileRepo: r,
		LabelRepo:   l,
		GroupRepo:   g,
		Validator:   v.NewNameValidator(),
		Auth:        auth.NewVoidAuth(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
