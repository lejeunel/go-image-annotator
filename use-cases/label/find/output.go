package find

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type OutputPort interface {
	SuccessFindLabel(lbl.Label)
	Error(error)
}
