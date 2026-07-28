package fake

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type ImageStore struct {
	Err       error
	Return    *im.Image
	DeletedId im.ImageId
}

func (s *ImageStore) Find(baseImage im.BaseImage) (*im.Image, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	if s.Return != nil {
		return s.Return, nil
	}
	return &im.Image{}, nil
}

func (s *ImageStore) DeleteAsset(id im.ImageId) error {
	s.DeletedId = id
	return nil
}
