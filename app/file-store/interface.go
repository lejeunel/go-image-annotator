package file_store

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"io"
)

type Interface interface {
	Store(im.ImageId, io.Reader) error
	Delete(im.ImageId) error
	Get(im.ImageId) (io.Reader, error)
}

type ReadInterface interface {
	Get(im.ImageId) (io.Reader, error)
}
