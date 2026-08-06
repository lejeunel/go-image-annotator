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
	MetaDataRepo
	KeyValidator   kv.Validator
	ValueValidator vv.Validator
	Auth
}

func New(c CollectionRepo, m MetaDataRepo,
	kv kv.Validator, vv vv.Validator,
	opts ...Option) Interactor {
	i := &Interactor{
		CollectionRepo: c,
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

	exists, err := i.MetaDataRepo.KeyExists(r.Collection, imageId, r.Key)
	if err != nil {
		out.Error(fmt.Errorf("%v: checking existence of key %v: %v: %w", errCtx, r.Key, err, e.ErrInternal))
		return
	}
	if *exists {
		out.Error(fmt.Errorf("%v: checking existence of key %v: %w", errCtx, r.Key, e.ErrValidation))
		return
	}

	if err := i.KeyValidator.Validate(r.Key); err != nil {
		out.Error(fmt.Errorf("%v: validating key %v: %w", errCtx, r.Key, err))
		return
	}
	if err := i.ValueValidator.Validate(r.Key); err != nil {
		out.Error(fmt.Errorf("%v: validating value %v: %w", errCtx, r.Key, err))
		return
	}
	if err := i.MetaDataRepo.Add(r.Collection, imageId, r.Key, r.Value); err != nil {
		out.Error(fmt.Errorf("%v: adding meta-data with key %v and value %v: %w",
			errCtx, r.Key, r.Value, err))
		return
	}

	out.SuccessAddMetadata()

}
