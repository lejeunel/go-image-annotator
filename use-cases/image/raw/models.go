package raw

import (
	"io"

	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Request struct {
	ImageId string
}

type Response struct {
	io.Reader
	im.Specs
}
