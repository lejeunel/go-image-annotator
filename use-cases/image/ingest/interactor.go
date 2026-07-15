package ingest

import (
	"context"
	"fmt"
	"io"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	ing "github.com/lejeunel/go-image-annotator/modules/ingester"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type IImageSpecsDetector interface {
	Detect(io.Reader) (*im.ImageSpecs, io.Reader, error)
}

type Ingester interface {
	Ingest() error
}

type Interactor struct {
	ingester ing.Interface
	auth     Auth
	repo     CollectionRepo
}
type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(ingester ing.Interface, repo CollectionRepo, opts ...Option) *Interactor {
	i := &Interactor{ingester: ingester,
		auth: auth.NewVoidAuth(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i Interactor) Execute(ctx context.Context, r ing.Request, out OutputPort) {
	errCtx := fmt.Errorf("ingesting image")
	collection, err := i.findCollectionByName(r.Collection)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if collection.Group != nil {
		if err := i.auth.IngestImage(ctx, collection.Group.Name); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	user := u.IdentityFromContext(ctx)
	if user == nil {
		out.Error(fmt.Errorf("%w: extracting user identity failed from context: %w", errCtx, e.ErrAuthentication))
		return
	}
	response, err := i.ingester.Ingest(ing.Request{UserId: user.Id, Collection: collection.Name, Labels: r.Labels,
		BoundingBoxes: r.BoundingBoxes, Reader: r.Reader})
	if err != nil {
		out.Error(fmt.Errorf("%w: %w", errCtx, err))
		return
	}

	out.Success(*response)

}

func (i Interactor) findCollectionByName(name string) (*clc.Collection, error) {
	collection, err := i.repo.FindCollectionByName(name)
	baseErr := fmt.Errorf("finding collection with name %v", name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return collection, nil

}
