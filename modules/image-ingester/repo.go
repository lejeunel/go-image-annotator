package ingester

import (
	"time"

	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type CollectionRepo interface {
	Find(string) (*clc.Collection, error)
}

type LabelRepo interface {
	FindLabel(string) (*lbl.Label, error)
}

type AnnotationRepo interface {
	AddImageLabel(im.ImageId, clc.CollectionName, an.ImageLabel, *u.UserId, *time.Time) error
	AddBoundingBox(im.ImageId, clc.CollectionName, an.BoundingBox, *u.UserId, *time.Time) error
	AddPolygon(im.ImageId, clc.CollectionName, an.Polygon, *u.UserId, *time.Time) error
}

type ImageRepo interface {
	AddImage(im.ImageId, []byte, im.Specs) error
	AddToCollection(im.ImageId, clc.CollectionName) error
	FindImageIdByHash([]byte) (*im.ImageId, error)
	Delete(im.ImageId) error
}
