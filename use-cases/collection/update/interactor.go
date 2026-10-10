package update

import (
	"context"
	"errors"
	"fmt"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	CollectionRepo
	GroupRepo
	ProfileRepo
	Auth
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(cr CollectionRepo, gr GroupRepo, pr ProfileRepo, opts ...Option) Interactor {
	i := &Interactor{cr, gr, pr, auth.NewVoidAuth()}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "updating collection"
	group, err := i.CollectionRepo.GetGroup(r.Name)
	if (err != nil) && !errors.Is(err, e.ErrNotFound) {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.Auth.UpdateCollection(ctx, group); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	response := clc.Collection{
		Name: r.NewName, Description: r.NewDescription,
		Profile: r.NewProfile, Group: r.NewGroup,
	}
	updateModel := clc.UpdateModel{
		Name: r.Name, NewName: r.NewName,
		NewDescription: r.NewDescription,
	}

	if err := i.ensureCollectionNameExists(r.Name); err != nil {
		fmt.Println(err)
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	updateModel.Name = r.Name

	if r.NewName != r.Name {
		if err := i.ensureCollectionNameDoesNotExist(r.NewName); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}
	updateModel.NewName = r.NewName

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
		if err := i.Auth.UpdateCollection(ctx, r.NewGroup); err != nil {
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

	if r.NewProfile != nil {
		profileExists, err := i.ProfileRepo.Exists(*r.NewProfile)
		if err != nil {
			out.Error(fmt.Errorf("%v: checking existence of profile: %w", errCtx, err))
			return
		}
		if !*profileExists {
			out.Error(
				fmt.Errorf(
					"%v: requested assignment to new profile %v: %w",
					errCtx,
					*r.NewGroup,
					err,
				),
			)
			return
		}
	}
	updateModel.NewProfile = r.NewProfile

	if err := i.CollectionRepo.Update(updateModel); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessUpdateCollection(response)
}

func (i Interactor) ensureCollectionNameExists(name string) error {
	baseErr := fmt.Errorf("ensuring that collection with name %v exists", name)
	exists, err := i.CollectionRepo.Exists(name)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, e.ErrInternal)
	}
	if !exists {
		return fmt.Errorf("%w: %w", baseErr, e.ErrNotFound)
	}
	return nil
}

func (i Interactor) ensureCollectionNameDoesNotExist(name string) error {
	baseErr := fmt.Errorf("ensuring that a collection with name %v does not already exist", name)
	exists, err := i.CollectionRepo.Exists(name)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, e.ErrInternal)
	}
	if exists {
		return fmt.Errorf("%w: %w", baseErr, e.ErrDuplicate)
	}
	return nil
}
