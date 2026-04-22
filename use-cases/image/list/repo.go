package list

import (
	ist "github.com/lejeunel/go-image-annotator-v2/app/image-store"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
)

type Repo interface {
	List(ist.FilteringParams) (*[]im.BaseImage, error)
	Count(ist.CountingParams) (*int64, error)
}
