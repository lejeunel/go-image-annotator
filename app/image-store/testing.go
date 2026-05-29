package image_store

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type FakeRepo struct {
	Err                    error
	ErrOnExists            bool
	Specs                  im.ImageSpecs
	ErrOnSpecs             bool
	ErrOnFindImageLabel    bool
	ErrOnFindBoundingBoxes bool
	MissingCollection      bool
	ErrOnFindCollection    bool
	Collection             clc.Collection
	Labels                 []*a.ImageLabel
	BoundingBoxes          []*a.BoundingBox
}

func (r *FakeRepo) ImageExistsInCollection(imageId im.ImageId, collectionId clc.CollectionId) (bool, error) {
	if r.ErrOnExists {
		return false, r.Err
	}
	return true, nil
}
func (r *FakeRepo) GetSpecs(imageId im.ImageId) (*im.ImageSpecs, error) {
	if r.ErrOnSpecs {
		return nil, r.Err
	}
	return &r.Specs, nil
}

func (r *FakeRepo) FindBoundingBoxes(imageId im.ImageId, collectionId clc.CollectionId) ([]*a.BoundingBox, error) {
	if r.ErrOnFindBoundingBoxes {
		return nil, r.Err
	}
	if r.BoundingBoxes != nil {
		return r.BoundingBoxes, nil
	}
	return nil, nil
}

func (r *FakeRepo) FindImageLabels(imageId im.ImageId, collectionId clc.CollectionId) ([]*a.ImageLabel, error) {
	if r.ErrOnFindImageLabel {
		return nil, r.Err
	}
	if r.Labels != nil {
		return r.Labels, nil
	}
	return nil, nil
}

func (r *FakeRepo) FindCollectionByName(name string) (*clc.Collection, error) {
	if r.MissingCollection {
		return nil, e.ErrNotFound
	}
	if r.ErrOnFindCollection {
		return nil, r.Err
	}
	return &r.Collection, nil
}

type FakeImageStore struct {
	Err    error
	Return *im.Image
}

func (s *FakeImageStore) Find(baseImage im.BaseImage) (*im.Image, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	if s.Return != nil {
		return s.Return, nil
	}
	return &im.Image{}, nil
}
