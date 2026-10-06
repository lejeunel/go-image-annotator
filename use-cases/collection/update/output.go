package update

import (
	c "github.com/lejeunel/go-image-annotator/entities/collection"
)

type OutputPort interface {
	SuccessUpdateCollection(c.Collection)
	Error(error)
}
