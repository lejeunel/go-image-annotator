package delete

import (
	"iter"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type ImageRepo interface {
	Iterate(im.FilterQueryStr, int) iter.Seq2[im.BaseImage, error]
}

type CollectionRepo interface {
	Find(string) (*clc.Collection, error)
	Delete(string) error
}
