package create

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type ProfileRepo interface {
	Create(pr.ProfileId, pr.ProfileName, pr.ProfileDescription) error
	Exists(pr.ProfileName) (*bool, error)
	AddLabel(pr.ProfileName, lbl.LabelName) error
}

type LabelRepo interface {
	Exists(lbl.LabelName) (bool, error)
}
