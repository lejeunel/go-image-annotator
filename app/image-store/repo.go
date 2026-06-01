package image_store

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Repo interface {
	FindCollectionByName(string) (*clc.Collection, error)
	FindImageLabels(im.ImageId, clc.CollectionId) ([]a.ImageLabel, error)
	FindBoundingBoxes(im.ImageId, clc.CollectionId) ([]a.BoundingBox, error)
	ImageExistsInCollection(im.ImageId, clc.CollectionId) (bool, error)
	GetSpecs(im.ImageId) (*im.ImageSpecs, error)
}
