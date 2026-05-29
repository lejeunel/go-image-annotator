package import_image

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Repo interface {
	ImageExists(im.ImageId) (bool, error)
	FindCollectionByName(string) (*clc.Collection, error)
	ImageExistsInCollection(im.ImageId, clc.CollectionId) (bool, error)
	AddToCollection(im.ImageId, clc.CollectionId) error
}
