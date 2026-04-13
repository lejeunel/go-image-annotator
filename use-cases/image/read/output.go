package read

import (
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
)

type OutputPort interface {
	SuccessReadImage(im.Image)
	Error(error)
}
