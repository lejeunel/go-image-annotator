package delete

import (
	"context"
	"errors"
	"fmt"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	sauth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Auth interface {
	DeleteMetadata(ctx context.Context, group string) error
}

type Interactor struct {
	CollectionRepo
	MetaDataRepo
	Auth
}

func New(c CollectionRepo, m MetaDataRepo, opts ...Option) Interactor {
	i := &Interactor{
		CollectionRepo: c,
		MetaDataRepo:   m,
		Auth:           sauth.NewVoidAuth(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "deleting metadata"
	group, err := i.CollectionRepo.GetGroup(r.Collection)
	if (err != nil) && !errors.Is(err, e.ErrNotFound) {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if group != nil {
		if err := i.Auth.DeleteMetadata(ctx, *group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: parsing image id %v: %w", errCtx, r.ImageId, err))
		return
	}

	exists, err := i.MetaDataRepo.KeyExists(r.Collection, imageId, r.Key)
	if err != nil {
		out.Error(fmt.Errorf("%v: checking existence of key %v: %w", errCtx, r.Key, err))
		return
	}
	if !exists {
		out.Error(fmt.Errorf("%v: checking existence of key %v: %w", errCtx, r.Key, e.ErrNotFound))
		return
	}

	if err := i.MetaDataRepo.Delete(r.Collection, imageId, r.Key); err != nil {
		out.Error(fmt.Errorf("%v: deleting key %v: %v: %w", errCtx, r.Key, err, e.ErrInternal))
		return

	}

	out.SuccessDeleteMetadata(r.Key)
}
