package fake

import (
	"slices"

	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type CreatedProfile struct {
	Id          pr.ProfileId
	Name        pr.ProfileName
	Description pr.ProfileDescription
}

type ProfileRepo struct {
	ErrOnCreate   error
	ErrOnExists   error
	ErrOnAddLabel error
	ErrOnFind     error
	ExistingNames []string
	Created       []CreatedProfile
	AddedLabels   []lbl.LabelName
	Return        pr.Profile
}

func (r *ProfileRepo) Create(
	id pr.ProfileId,
	name pr.ProfileName,
	desc pr.ProfileDescription,
) error {
	if r.ErrOnCreate != nil {
		return r.ErrOnCreate
	}

	r.Created = append(r.Created,
		CreatedProfile{Id: id, Name: name, Description: desc})
	return nil
}

func (r *ProfileRepo) Exists(name string) (*bool, error) {
	if r.ErrOnExists != nil {
		return nil, r.ErrOnExists
	}

	exist := true
	if slices.Contains(r.ExistingNames, name) {
		return &exist, nil
	}
	exist = false
	return &exist, nil
}

func (r *ProfileRepo) AddLabel(profileName pr.ProfileName, label lbl.LabelName) error {
	if r.ErrOnAddLabel != nil {
		return r.ErrOnAddLabel
	}
	r.AddedLabels = append(r.AddedLabels, label)
	return nil
}

func (r *ProfileRepo) Find(name string) (*pr.Profile, error) {
	if r.ErrOnFind != nil {
		return nil, r.ErrOnFind
	}

	return &r.Return, nil
}
