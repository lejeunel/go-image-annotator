package find

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type OutputPort interface {
	SuccessReadImage(im.Image)
	Error(error)
}
