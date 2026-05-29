package list

import (
	ist "github.com/lejeunel/go-image-annotator/app/image-store"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Repo interface {
	List(ist.FilteringParams) (*[]im.BaseImage, error)
	Count(ist.CountingParams) (*int64, error)
}
