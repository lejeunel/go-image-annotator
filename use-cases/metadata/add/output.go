package add

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type OutputPort interface {
	Error(error)
	SuccessAddMetadata(im.Image)
}
