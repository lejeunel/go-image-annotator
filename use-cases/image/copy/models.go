package copy

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Request struct {
	ImageId     string
	Source      string
	Destination string
	Deep        bool
}

type Response struct {
	Image im.Image
}
