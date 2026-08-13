package list

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Repo interface {
	Slice(im.FilterQueryStr, pa.PaginationParams, im.OrderingStr) ([]im.BaseImage, error)
	Count(im.FilterQueryStr) (*int64, error)
}
