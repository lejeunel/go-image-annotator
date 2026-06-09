package create

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
)

type GroupRepo interface {
	Find(string) (*grp.Group, error)
}

type CollectionRepo interface {
	Create(clc.Collection) error
	Exists(string) (bool, error)
}
