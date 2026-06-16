package update_label

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"time"
)

type Repo interface {
	FindLabel(string) (*lbl.Label, error)
	UpdateLabelOfAnnotation(a.AnnotationId, lbl.LabelId, *u.UserId, *time.Time) error
	GroupOfAnnotation(a.AnnotationId) (*string, error)
}
