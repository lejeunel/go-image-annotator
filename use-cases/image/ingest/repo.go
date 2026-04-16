package ingest

import (
	an "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator-v2/entities/collection"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	lbl "github.com/lejeunel/go-image-annotator-v2/entities/label"
)

type CollectionRepo interface {
	FindCollectionByName(string) (*clc.Collection, error)
}

type LabelRepo interface {
	FindLabelByName(string) (*lbl.Label, error)
}

type AnnotationRepo interface {
	AddImageLabel(an.AnnotationId, im.ImageId, clc.CollectionId, lbl.LabelId) error
	AddBoundingBox(im.ImageId, clc.CollectionId, an.BoundingBox) error
}

type ImageRepo interface {
	AddImage(im.ImageId, []byte, im.ImageSpecs) error
	AddToCollection(im.ImageId, clc.CollectionId) error
	FindImageIdByHash([]byte) (*im.ImageId, error)
	Delete(im.ImageId) error
}
