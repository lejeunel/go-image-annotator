package create

import (
	c "github.com/lejeunel/go-image-annotator/entities/collection"
)

type OutputPort interface {
	SuccessCreateCollection(c.Collection)
	Error(error)
}
