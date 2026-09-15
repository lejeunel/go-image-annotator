package update

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type ProfileRepo interface {
	Update(pr.UpdateModel) error
	Exists(string) (*bool, error)
	GetGroup(string) (*string, error)
}

type LabelRepo interface {
	Exists(string) (bool, error)
}

type GroupRepo interface {
	Exists(string) (*bool, error)
}
