package import_image

import (
	"context"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success(Response) {
	p.GotSuccess = true
}

type FakeRepo struct {
	ErrOnImageExists             bool
	ErrOnFindCollection          bool
	ErrOnImageExistsInCollection bool
	ErrOnImport                  bool
	ImageMissing                 bool
	ImageAlreadyInCollection     bool

	ImportedImageId          im.ImageId
	ImportedIntoCollectionId clc.CollectionId

	Err                   error
	DestinationCollection clc.Collection
}

func (r *FakeRepo) ImageExists(imageId im.ImageId) (bool, error) {
	if r.ErrOnImageExists {
		return false, r.Err
	}
	if r.ImageMissing {
		return false, nil
	}
	return true, nil
}

func (r *FakeRepo) FindCollectionByName(name string) (*clc.Collection, error) {
	if r.ErrOnFindCollection {
		return nil, r.Err
	}
	return &r.DestinationCollection, nil
}

func (r *FakeRepo) ImageExistsInCollection(imageId im.ImageId, collectionId clc.CollectionId) (bool, error) {
	if r.ErrOnImageExistsInCollection {
		return false, r.Err
	}
	if r.ImageAlreadyInCollection {
		return true, nil
	}
	return false, nil
}

func (r *FakeRepo) AddToCollection(imageId im.ImageId, collectionId clc.CollectionId) error {
	if r.ErrOnImport {
		return r.Err
	}
	r.ImportedImageId = imageId
	r.ImportedIntoCollectionId = collectionId
	return nil
}

type FailingAuth struct {
}

func (f FailingAuth) ImportImage(ctx context.Context, srcGroup string, dstGroup string) error {
	return e.ErrAuth
}
