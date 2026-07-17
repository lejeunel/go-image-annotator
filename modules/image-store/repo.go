package image_store

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type AnnotationRepo interface {
	FindImageLabels(im.ImageId, clc.CollectionId) ([]a.ImageLabel, error)
	FindBoundingBoxes(im.ImageId, clc.CollectionId) ([]a.BoundingBox, error)
	FindPolygons(im.ImageId, clc.CollectionId) ([]a.Polygon, error)
}

type CollectionRepo interface {
	FindCollectionByName(string) (*clc.Collection, error)
}

type ImageRepo interface {
	GetSpecs(im.ImageId) (*im.ImageSpecs, error)
	ImageExistsInCollection(im.ImageId, clc.CollectionId) (bool, error)
}
