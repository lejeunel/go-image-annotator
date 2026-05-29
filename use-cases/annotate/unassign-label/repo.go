package unassign_label

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type Repo interface {
	RemoveImageLabel(im.ImageId, clc.CollectionId, lbl.LabelId) error
}
