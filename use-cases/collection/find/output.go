package find

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
)

type OutputPort interface {
	Error(error)
	SuccessFindCollection(clc.Collection)
}
