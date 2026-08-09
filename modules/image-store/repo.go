package image_store

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
)

type AnnotationRepo interface {
	FindImageLabels(im.ImageId, clc.CollectionId) ([]a.ImageLabel, error)
	FindBoundingBoxes(im.ImageId, clc.CollectionId) ([]a.BoundingBox, error)
	FindPolygons(im.ImageId, clc.CollectionId) ([]a.Polygon, error)
	RemoveAllAnnotations(im.ImageId, clc.CollectionName) error
}

type CollectionRepo interface {
	Find(string) (*clc.Collection, error)
}

type ImageRepo interface {
	GetSpecs(im.ImageId) (*im.Specs, error)
	ImageExistsInCollection(im.ImageId, clc.CollectionName) (bool, error)
	RemoveImageFromCollection(im.ImageId, clc.CollectionName) error
	IsUsed(im.ImageId) (*bool, error)
}

type MetaRepo interface {
	List(clc.CollectionName, im.ImageId) ([]m.MetaData, error)
}
