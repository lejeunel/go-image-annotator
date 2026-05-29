package update_label

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type Repo interface {
	FindLabelByName(string) (*lbl.Label, error)
	UpdateLabelOfAnnotation(a.AnnotationId, lbl.LabelId) error
}
