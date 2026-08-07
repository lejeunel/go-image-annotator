package update

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	sauth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Auth interface {
	UpdateMetadata(ctx context.Context, group string) error
}

type Interactor struct {
	CollectionRepo
	ImageRepo
	MetaDataRepo
	Auth
}

func New(c CollectionRepo, ir ImageRepo, m MetaDataRepo,
	opts ...Option) Interactor {
	i := &Interactor{
		CollectionRepo: c,
		MetaDataRepo:   m,
		Auth:           sauth.NewVoidAuth()}
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

	errCtx := "updating metadata"
	group, err := i.CollectionRepo.GetGroup(r.Collection)
	if (err != nil) && !(errors.Is(err, e.ErrNotFound)) {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if group != nil {
		if err := i.Auth.UpdateMetadata(ctx, *group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: parsing image id %v: %w", errCtx, imageId, err))
		return
	}

	collectionExists, err := i.CollectionRepo.Exists(r.Collection)
	if err != nil {
		out.Error(fmt.Errorf("%v: checking existence of collection %v: %v: %w", errCtx, r.Collection, err, e.ErrInternal))
		return
	}
	if !collectionExists {
		out.Error(fmt.Errorf("%v: checking existence of collection %v: %v: %w", errCtx, r.Collection, err, e.ErrValidation))
		return
	}

	imageInCollection, err := i.ImageRepo.ImageExistsInCollection(imageId, r.Collection)
	if err != nil {
		out.Error(fmt.Errorf("%v: checking whether image %v is in collection %v: %v: %w", errCtx, imageId, r.Collection, err, e.ErrInternal))
		return
	}
	if !imageInCollection {
		out.Error(fmt.Errorf("%v: checking whether image %v is in collection %v: %v: %w", errCtx, imageId, r.Collection, err, e.ErrValidation))
		return
	}

	exists, err := i.MetaDataRepo.KeyExists(r.Collection, imageId, r.Key)
	if err != nil {
		out.Error(fmt.Errorf("%v: checking existence of key %v: %v: %w", errCtx, r.Key, err, e.ErrInternal))
		return
	}
	if !*exists {
		out.Error(fmt.Errorf("%v: checking existence of key %v: %w", errCtx, r.Key, e.ErrValidation))
		return
	}

	currentValue, err := i.MetaDataRepo.GetValue(r.Collection, imageId, r.Key)
	if err != nil {
		out.Error(fmt.Errorf("%v: fetching current value at key %v: %v: %w", errCtx, r.Key, err, e.ErrInternal))
		return
	}

	if reflect.TypeOf(*currentValue) != reflect.TypeOf(r.Value) {
		out.Error(fmt.Errorf("%v: comparing types of current value %v with new value %v: %v: %w",
			errCtx, *currentValue, r.Value, err, e.ErrValidation))
		return
	}

	if err := i.MetaDataRepo.UpdateValue(r.Collection, imageId, r.Key, r.Value); err != nil {
		out.Error(fmt.Errorf("%v: updating key %v with new value %v: %v: %w",
			errCtx, r.Key, r.Value, err, e.ErrInternal))
		return
	}

	out.SuccessUpdateMetadata()

}
