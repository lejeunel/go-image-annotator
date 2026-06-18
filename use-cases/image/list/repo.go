package list

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Repo interface {
	List(im.FilteringParams, im.OrderingParams) ([]im.BaseImage, error)
	Count(im.CountingParams) (*int64, error)
}
