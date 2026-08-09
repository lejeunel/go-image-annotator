package ingest

import (
	"context"
	"fmt"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	ing "github.com/lejeunel/go-image-annotator/modules/image-ingester"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Ingester interface {
	Ingest(ing.Request) (*ing.Response, error)
}

type Interactor struct {
	Ingester
	auth Auth
	repo CollectionRepo
}
type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(ingester Ingester, repo CollectionRepo, opts ...Option) *Interactor {
	i := &Interactor{
		Ingester: ingester,
		repo:     repo,
		auth:     auth.NewVoidAuth(),
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
		if err := i.auth.IngestImage(ctx, *collection.Group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	user := u.IdentityFromContext(ctx)
	if user == nil {
		out.Error(
			fmt.Errorf(
				"%w: extracting user identity failed from context: %w",
				errCtx,
				e.ErrAuthentication,
			),
		)
		return
	}
	response, err := i.Ingester.Ingest(ing.Request{
		UserId: user.Id, Collection: collection.Name, Labels: r.Labels,
		BoundingBoxes: r.BoundingBoxes, Reader: r.Reader,
	})
	if err != nil {
		out.Error(fmt.Errorf("%w: %w", errCtx, err))
		return
	}

	out.Success(*response)
}

func (i Interactor) findCollectionByName(name string) (*clc.Collection, error) {
	collection, err := i.repo.Find(name)
	baseErr := fmt.Errorf("finding collection with name %v", name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return collection, nil
}
