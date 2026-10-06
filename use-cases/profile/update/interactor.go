package update

import (
	"context"
	"errors"
	"fmt"

	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	ProfileRepo
	GroupRepo
	LabelRepo
	Auth
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(pr ProfileRepo, gr GroupRepo, lr LabelRepo, opts ...Option) Interactor {
	i := &Interactor{pr, gr, lr, auth.NewVoidAuth()}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "updating profile"
	src, err := i.ProfileRepo.Find(r.Name)
	if err != nil {
		out.Error(
			fmt.Errorf("%v: fetching profile %v: %w", errCtx, r.Name, e.ErrValidation),
		)
		return
	}

	if src.Group != nil {
		if err := i.Auth.UpdateProfile(ctx, src.Group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}
	updateModel := pr.UpdateModel{Name: r.Name, NewDescription: r.NewDescription}

	if r.NewName != r.Name {
		_, err := i.ProfileRepo.Find(r.NewName)
		if !errors.Is(err, e.ErrNotFound) {
			out.Error(
				fmt.Errorf(
					"%v: making sure that profile %v does not exist: %w",
					errCtx,
					r.NewName,
					e.ErrValidation,
				),
			)
			return
		}
	}

	updateModel.NewName = r.NewName

	for _, label := range r.NewLabels {
		labelExists, err := i.LabelRepo.Exists(label)
		if err != nil {
			out.Error(fmt.Errorf("%v: checking existence of label %v: %w", errCtx, label, err))
			return
		}
		if !labelExists {
			out.Error(
				fmt.Errorf(
					"%v: checking existence of label %v: %w",
					errCtx,
					label,
					e.ErrValidation,
				),
			)
			return
		}

	}
	updateModel.NewLabels = r.NewLabels

	if r.NewGroup != nil {
		groupExists, err := i.GroupRepo.Exists(*r.NewGroup)
		if err != nil {
			out.Error(fmt.Errorf("%v: checking existence of group: %w", errCtx, err))
			return
		}
		if !*groupExists {
			out.Error(
				fmt.Errorf(
					"%v: requested assignment to new group %v: %w",
					errCtx,
					*r.NewGroup,
					err,
				),
			)
			return
		}
		if err := i.Auth.UpdateProfile(ctx, r.NewGroup); err != nil {
			out.Error(
				fmt.Errorf(
					"%v: authorizing assignment to new group %v: %w",
					errCtx,
					*r.NewGroup,
					err,
				),
			)
			return
		}
	}
	updateModel.NewGroup = r.NewGroup

	if err := i.ProfileRepo.Update(updateModel); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessUpdateProfile(pr.Profile{
		Id: src.Id, Name: r.NewName, Description: r.NewDescription,
		Group: r.NewGroup, Labels: r.NewLabels,
	})
}
