package clone

import (
	"iter"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"time"
)

type ImageRepo interface {
	Iterate(im.FilteringParams, int) iter.Seq2[im.BaseImage, error]
	AddToCollection(im.ImageId, clc.CollectionId) error
}
type CollectionRepo interface {
	Create(clc.Collection) error
	Exists(string) (bool, error)
}

type AnnotationRepo interface {
	AddImageLabel(im.ImageId, clc.CollectionId, a.ImageLabel, *u.UserId, *time.Time) error
	AddBoundingBox(im.ImageId, clc.CollectionId, a.BoundingBox, *u.UserId, *time.Time) error
	AddPolygon(im.ImageId, clc.CollectionId, a.Polygon, *u.UserId, *time.Time) error
}

type GroupRepo interface {
	Find(string) (*grp.Group, error)
}
