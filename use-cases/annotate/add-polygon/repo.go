package add_polygon

import (
	"time"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	AddPolygon(im.ImageId, clc.CollectionId, a.Polygon, *u.UserId, *time.Time) error
}

type LabelRepo interface {
	FindLabel(string) (*lbl.Label, error)
}
