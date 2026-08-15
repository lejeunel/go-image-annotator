package list

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Repo interface {
	Slice(im.FilterStr, pa.PaginationParams, im.OrderStr) ([]im.BaseImage, error)
	Count(im.FilterStr) (*int64, error)
}
