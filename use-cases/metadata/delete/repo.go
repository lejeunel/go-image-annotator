package delete

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type CollectionRepo interface {
	GetGroup(string) (*string, error)
}

type MetaDataRepo interface {
	Delete(clc.CollectionName, im.ImageId, string) error
	KeyExists(clc.CollectionName, im.ImageId, string) (*bool, error)
}
