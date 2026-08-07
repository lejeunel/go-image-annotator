package modify_polygon

import (
	"time"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type AnnotationRepo interface {
	UpdatePolygon(a.AnnotationId, a.PolygonUpdatables, *u.UserId, *time.Time) error
	GroupOfAnnotation(a.AnnotationId) (*string, error)
}

type LabelRepo interface {
	FindLabel(string) (*lbl.Label, error)
}
