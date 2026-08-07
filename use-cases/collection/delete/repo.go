package delete

import (
	"iter"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Repos struct {
	ImageRepo
	CollectionRepo
	AnnotationRepo
}

type ImageRepo interface {
	Iterate(im.Filtering, int) iter.Seq2[im.BaseImage, error]
	RemoveImageFromCollection(im.ImageId, clc.CollectionId) error
	IsUsed(im.ImageId) (*bool, error)
}

type CollectionRepo interface {
	Find(string) (*clc.Collection, error)
	Delete(string) error
	IsPopulated(string) (*bool, error)
}

type AnnotationRepo interface {
	RemoveAllAnnotations(im.ImageId, string) error
}
