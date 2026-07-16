package import_image

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type ImageRepo interface {
	ImageExists(im.ImageId) (bool, error)
	ImageExistsInCollection(im.ImageId, clc.CollectionId) (bool, error)
	AddToCollection(im.ImageId, clc.CollectionId) error
}

type CollectionRepo interface {
	FindCollectionByName(string) (*clc.Collection, error)
}
