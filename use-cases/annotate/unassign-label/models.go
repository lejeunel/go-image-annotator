package unassign_label

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Response struct{}

type Request struct {
	ImageId    im.ImageId
	Collection string
	Label      string
}
