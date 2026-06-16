package assign_label

import (
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"time"
)

type AnnotationRepo interface {
	AddImageLabel(im.ImageId, clc.CollectionId, an.ImageLabel, *u.UserId, *time.Time) error
}

type LabelRepo interface {
	FindLabel(string) (*lbl.Label, error)
}
