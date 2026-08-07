package add

import (
	"context"
	"errors"
	"fmt"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	sauth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	kv "github.com/lejeunel/go-image-annotator/modules/string-validator"
	vv "github.com/lejeunel/go-image-annotator/modules/value-validator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Auth interface {
	AddMetadata(ctx context.Context, group string) error
}

type Interactor struct {
	CollectionRepo
	ImageRepo
	MetaDataRepo
	KeyValidator   kv.Validator
	ValueValidator vv.Validator
	Auth
}

func New(c CollectionRepo, ir ImageRepo,
	m MetaDataRepo,
	kv kv.Validator, vv vv.Validator,
	opts ...Option) Interactor {
	i := &Interactor{
		CollectionRepo: c,
		ImageRepo:      ir,
		MetaDataRepo:   m,
		KeyValidator:   kv,
		ValueValidator: vv,

		Auth: sauth.NewVoidAuth()}
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

	errCtx := "adding metadata"
	group, err := i.CollectionRepo.GetGroup(r.Collection)
	if (err != nil) && !(errors.Is(err, e.ErrNotFound)) {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if group != nil {
		if err := i.Auth.AddMetadata(ctx, *group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: parsing image id %v: %w", errCtx, imageId, err))
		return
	}

	keyExists, err := i.MetaDataRepo.KeyExists(r.Collection, imageId, r.Key)
	if err != nil {
		out.Error(fmt.Errorf("%v: checking existence of key %v: %v: %w", errCtx, r.Key, err, e.ErrInternal))
		return
	}
	if keyExists {
		out.Error(fmt.Errorf("%v: checking existence of key %v: %w", errCtx, r.Key, e.ErrValidation))
		return
	}

	collectionExists, err := i.CollectionRepo.Exists(r.Collection)
	if err != nil {
		out.Error(fmt.Errorf("%v: checking existence of collection %v: %v: %w", errCtx, r.Collection, err, e.ErrInternal))
		return
	}
	if !collectionExists {
		out.Error(fmt.Errorf("%v: checking existence of collection %v: %w", errCtx, r.Collection, e.ErrValidation))
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

	if err := i.KeyValidator.Validate(r.Key); err != nil {
		out.Error(fmt.Errorf("%v: validating key %v: %w", errCtx, r.Key, err))
		return
	}
	if err := i.ValueValidator.Validate(r.Value); err != nil {
		out.Error(fmt.Errorf("%v: validating value %v: %w", errCtx, r.Value, err))
		return
	}
	if err := i.MetaDataRepo.Add(r.Collection, imageId, r.Key, r.Value); err != nil {
		out.Error(fmt.Errorf("%v: adding meta-data with key %v and value %v: %w",
			errCtx, r.Key, r.Value, err))
		return
	}

	out.SuccessAddMetadata()

}
