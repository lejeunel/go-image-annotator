package list

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
)

type MetaDataRepo interface {
	List(clc.CollectionName, im.ImageId) ([]m.MetaData, error)
}
