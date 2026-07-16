package delete

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type ImageRepo interface {
	RemoveImageFromCollection(im.ImageId, clc.CollectionId) error
}
type AnnotationRepo interface {
	RemoveAnnotation(a.AnnotationId) error
}
