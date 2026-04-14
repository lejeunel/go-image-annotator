package add_bbox

import (
	a "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
)

type OutputPort interface {
	Error(error)
	SuccessAddBox(a.BoundingBox)
}
