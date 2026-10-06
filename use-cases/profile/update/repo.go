package update

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type ProfileRepo interface {
	Update(pr.UpdateModel) error
	Find(string) (*pr.Profile, error)
	GetGroup(string) (*string, error)
}

type LabelRepo interface {
	Exists(string) (bool, error)
}

type GroupRepo interface {
	Exists(string) (*bool, error)
}
