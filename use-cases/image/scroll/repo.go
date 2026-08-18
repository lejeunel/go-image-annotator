package scroll

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type ImageRepo interface {
	GetAdjacent(im.ImageId, clc.CollectionName, im.FilterStr, im.OrderStr, im.ScrollingDirection) (*im.BaseImage, error)
	ImageExists(imageId im.ImageId) (bool, error)
}
