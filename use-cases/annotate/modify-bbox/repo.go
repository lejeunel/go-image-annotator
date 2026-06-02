package modify_bbox

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type Repo interface {
	FindLabel(string) (*lbl.Label, error)
	UpdateBoundingBox(a.AnnotationId, a.BoundingBoxUpdatables) error
	GroupOfAnnotation(a.AnnotationId) (*string, error)
}
