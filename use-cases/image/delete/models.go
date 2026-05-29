package delete

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Request struct {
	ImageId    im.ImageId
	Collection string
}

type Response struct{}
