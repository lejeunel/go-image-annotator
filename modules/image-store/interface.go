package image_store

import im "github.com/lejeunel/go-image-annotator/entities/image"

type Interface interface {
	Find(base im.BaseImage) (*im.Image, error)
}
