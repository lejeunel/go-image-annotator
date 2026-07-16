package update

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
)

type GroupRepo interface {
	GroupOfCollection(string) (*string, error)
}

type CollectionRepo interface {
	Update(clc.UpdateModel) error
	Exists(string) (bool, error)
}
