package list

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Repo interface {
	Slice(im.FilteringParams, pa.PaginationParams, im.OrderingParams) ([]im.BaseImage, error)
	Count(im.CountingParams) (*int64, error)
}
