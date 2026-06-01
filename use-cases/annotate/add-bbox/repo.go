package add_bbox

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type Repo interface {
	AddBoundingBox(im.ImageId, clc.CollectionId, a.BoundingBox) error
	FindLabel(string) (*lbl.Label, error)
}
