package add_bbox

import (
	"time"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	AddBoundingBox(im.ImageId, clc.CollectionId, a.BoundingBox, *u.UserId, *time.Time) error
	FindLabel(string) (*lbl.Label, error)
}
