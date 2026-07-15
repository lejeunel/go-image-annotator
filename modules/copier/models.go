package copier

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Request struct {
	ImageId               im.ImageId
	SourceCollection      string
	DestinationCollection string
	Deep                  bool
}

type Response struct {
	ImageId    im.ImageId
	Collection string
}
