package create

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type OutputPort interface {
	SuccessCreateLabel(lbl.Label)
	Error(error)
}
