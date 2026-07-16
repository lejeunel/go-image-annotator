package fake

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type ImageRepo struct {
	Err              error
	RemovedImageId   im.ImageId
	ErrOnRemoveImage error
}

func (r *ImageRepo) RemoveImageFromCollection(imageId im.ImageId, collectionId clc.CollectionId) error {
	if r.ErrOnRemoveImage != nil {
		return r.ErrOnRemoveImage
	}
	r.RemovedImageId = imageId
	return nil
}
