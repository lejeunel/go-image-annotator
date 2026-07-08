package raw

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"io"
)

type Request struct {
	ImageId string
}

type Response struct {
	io.Reader
	im.ImageSpecs
}
