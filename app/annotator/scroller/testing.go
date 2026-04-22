package scroller

import (
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
)

type FakeRepo struct {
	Err                   error
	ErrOnImageExists      bool
	ErrOnCollectionExists bool
	NextImage             *im.BaseImage
	PreviousImage         *im.BaseImage
}

func (r *FakeRepo) ImageMustExist(id im.ImageId) error {
	if r.ErrOnImageExists {
		return r.Err
	}
	return nil
}

func (r *FakeRepo) CollectionMustExist(collection string) error {
	if r.ErrOnCollectionExists {
		return r.Err
	}
	return nil
}

func (r *FakeRepo) GetAdjacent(id im.ImageId, criteria ScrollingCriteria, d ScrollingDirection) (*im.BaseImage, error) {
	if d == ScrollNext {
		return r.NextImage, nil
	}
	if d == ScrollPrevious {
		return r.PreviousImage, nil
	}
	return nil, nil
}
