package scroller

import (
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
)

type Repo interface {
	GetAdjacent(im.ImageId, ScrollingCriteria, ScrollingDirection) (*im.BaseImage, error)
	ImageMustExist(imageId im.ImageId) error
	CollectionMustExist(string) error
}
