package add

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
)

type ImageRepo interface {
	ImageExistsInCollection(im.ImageId, clc.CollectionName) (bool, error)
}

type CollectionRepo interface {
	GetGroup(string) (*string, error)
	Exists(string) (bool, error)
}

type MetaDataRepo interface {
	Add(clc.CollectionName, im.ImageId, string, any) error
	List(clc.CollectionName, im.ImageId) ([]m.MetaData, error)
}
