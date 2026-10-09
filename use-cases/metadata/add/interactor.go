package add

import (
	"context"
	"fmt"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	sauth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	kv "github.com/lejeunel/go-image-annotator/modules/string-validator"
	vv "github.com/lejeunel/go-image-annotator/modules/value-validator"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Auth interface {
	AddMetadata(ctx context.Context, group *string) error
}

type ImageStore interface {
	Find(im.BaseImage) (*im.Image, error)
}

type Interactor struct {
	ImageStore
	MetaDataRepo
	KeyValidator   kv.Validator
	ValueValidator vv.Validator
	Auth
}

func New(s ImageStore,
	m MetaDataRepo,
	kv kv.Validator, vv vv.Validator,
	opts ...Option,
) Interactor {
	i := &Interactor{
		ImageStore:     s,
		MetaDataRepo:   m,
		KeyValidator:   kv,
		ValueValidator: vv,

		Auth: sauth.NewVoidAuth(),
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
	errCtx := "adding metadata"
	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	image, err := i.ImageStore.Find(im.BaseImage{})
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.Auth.AddMetadata(ctx, image.Collection.Group); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
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
	meta, err := i.MetaDataRepo.List(r.Collection, imageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: fetching updated meta-data: %w",
			errCtx, err))
		return
	}

	image.Meta = meta

	out.SuccessAddMetadata(*image)
}
