package ingest

import (
	"time"

	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type CollectionRepo interface {
	FindCollectionByName(string) (*clc.Collection, error)
}

type LabelRepo interface {
	FindLabel(string) (*lbl.Label, error)
}

type AnnotationRepo interface {
	AddImageLabel(im.ImageId, clc.CollectionId, an.ImageLabel, *u.UserId, *time.Time) error
	AddBoundingBox(im.ImageId, clc.CollectionId, an.BoundingBox, *u.UserId, *time.Time) error
}

type ImageRepo interface {
	AddImage(im.ImageId, []byte, im.ImageSpecs) error
	AddToCollection(im.ImageId, clc.CollectionId) error
	FindImageIdByHash([]byte) (*im.ImageId, error)
	Delete(im.ImageId) error
}
