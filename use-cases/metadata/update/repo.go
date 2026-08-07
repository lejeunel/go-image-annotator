package update

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type ImageRepo interface {
	ImageExistsInCollection(im.ImageId, clc.CollectionName) (bool, error)
}

type CollectionRepo interface {
	GetGroup(string) (*string, error)
	Exists(string) (bool, error)
}

type MetaDataRepo interface {
	GetValue(clc.CollectionName, im.ImageId, string) (*any, error)
	KeyExists(clc.CollectionName, im.ImageId, string) (*bool, error)
	UpdateValue(clc.CollectionName, im.ImageId, string, any) error
}
