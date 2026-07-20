package list

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Repo interface {
	Slice(im.Filtering, pa.PaginationParams, im.Ordering) ([]im.BaseImage, error)
	Count(im.Filtering) (*int64, error)
}
