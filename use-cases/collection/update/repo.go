package update

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
)

type CollectionRepo interface {
	Update(clc.UpdateModel) error
	Exists(string) (bool, error)
	GetGroup(string) (*string, error)
}

type GroupRepo interface {
	Exists(string) (*bool, error)
}
