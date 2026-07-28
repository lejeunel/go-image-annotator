package delete

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"iter"
)

type ImageRepo interface {
	Iterate(im.Filtering, int) iter.Seq2[im.BaseImage, error]
	RemoveImageFromCollection(im.ImageId, clc.CollectionId) error
}

type CollectionRepo interface {
	Find(string) (*clc.Collection, error)
	Delete(string) error
	IsPopulated(string) (*bool, error)
}

type GroupRepo interface {
	GroupOfCollection(string) (*string, error)
}

type AnnotationRepo interface {
	RemoveAllAnnotations(im.ImageId, string) error
}
