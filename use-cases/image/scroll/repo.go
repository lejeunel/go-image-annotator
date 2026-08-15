package scroll

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type ImageRepo interface {
	GetAdjacent(im.ImageId, im.FilterStr, im.OrderStr, im.ScrollingDirection) (*im.BaseImage, error)
	ImageExists(imageId im.ImageId) (bool, error)
}
