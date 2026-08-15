package clone

import (
	"iter"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type ImageRepo interface {
	Iterate(im.FilterStr, int) iter.Seq2[im.BaseImage, error]
	AddToCollection(im.ImageId, clc.CollectionName) error
}
type CollectionRepo interface {
	Create(clc.Collection) error
	Exists(string) (bool, error)
}

type GroupRepo interface {
	Find(string) (*grp.Group, error)
}
