package modify_bbox

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type Repo interface {
	FindLabelByName(string) (*lbl.Label, error)
	UpdateBoundingBox(a.AnnotationId, a.BoundingBoxUpdatables) error
}
