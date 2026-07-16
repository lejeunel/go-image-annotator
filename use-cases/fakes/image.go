package fake

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type ImageRepo struct {
	Err                          error
	RemovedImageId               im.ImageId
	ErrOnRemoveImage             error
	ErrOnImageExistsInCollection error
	ErrOnImageExists             error
	ErrOnImport                  error
	ImportedImageId              im.ImageId
	ImportedIntoCollectionId     clc.CollectionId
	ImageAlreadyInCollection     bool
}

func (r *ImageRepo) RemoveImageFromCollection(imageId im.ImageId, collectionId clc.CollectionId) error {
	if r.ErrOnRemoveImage != nil {
		return r.ErrOnRemoveImage
	}
	r.RemovedImageId = imageId
	return nil
}

func (r *ImageRepo) ImageExists(imageId im.ImageId) (bool, error) {
	if r.ErrOnImageExists != nil {
		return false, r.ErrOnImageExists
	}
	return true, nil
}
func (r *ImageRepo) ImageExistsInCollection(imageId im.ImageId, collectionId clc.CollectionId) (bool, error) {
	if r.ErrOnImageExistsInCollection != nil {
		return false, r.ErrOnImageExistsInCollection
	}
	if r.ImageAlreadyInCollection {
		return true, nil
	}
	return false, nil
}

func (r *ImageRepo) AddToCollection(imageId im.ImageId, collectionId clc.CollectionId) error {
	if r.ErrOnImport != nil {
		return r.ErrOnImport
	}
	r.ImportedImageId = imageId
	r.ImportedIntoCollectionId = collectionId
	return nil
}
