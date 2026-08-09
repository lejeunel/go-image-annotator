package image_store

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"time"
)

type AnnotationRepo interface {
	FindImageLabels(im.ImageId, clc.CollectionName) ([]a.ImageLabel, error)
	FindBoundingBoxes(im.ImageId, clc.CollectionName) ([]a.BoundingBox, error)
	FindPolygons(im.ImageId, clc.CollectionName) ([]a.Polygon, error)
	RemoveAllAnnotations(im.ImageId, clc.CollectionName) error
	AddImageLabel(im.ImageId, clc.CollectionName, a.ImageLabel, *u.UserId, *time.Time) error
	AddBoundingBox(im.ImageId, clc.CollectionName, a.BoundingBox, *u.UserId, *time.Time) error
	AddPolygon(im.ImageId, clc.CollectionName, a.Polygon, *u.UserId, *time.Time) error
}

type CollectionRepo interface {
	Find(string) (*clc.Collection, error)
}

type ImageRepo interface {
	GetSpecs(im.ImageId) (*im.Specs, error)
	ImageExistsInCollection(im.ImageId, clc.CollectionName) (bool, error)
	RemoveImageFromCollection(im.ImageId, clc.CollectionName) error
	IsUsed(im.ImageId) (*bool, error)
	AddToCollection(im.ImageId, clc.CollectionName) error
}

type MetaRepo interface {
	List(clc.CollectionName, im.ImageId) ([]m.MetaData, error)
	DeleteAll(clc.CollectionName, im.ImageId) error
}
