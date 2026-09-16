package create

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type ProfileRepo interface {
	Create(pr.Profile) error
	Exists(pr.ProfileName) (*bool, error)
}

type LabelRepo interface {
	Exists(lbl.LabelName) (bool, error)
}

type GroupRepo interface {
	Exists(string) (*bool, error)
}
