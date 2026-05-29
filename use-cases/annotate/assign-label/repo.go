package assign_label

import (
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type Repo interface {
	AddImageLabel(im.ImageId, clc.CollectionId, an.ImageLabel) error
	FindLabel(string) (*lbl.Label, error)
}
