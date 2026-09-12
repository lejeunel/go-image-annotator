package fake

import (
	"slices"

	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type ProfileRepo struct {
	ErrOnCreate   error
	ErrOnExists   error
	ExistingNames []string
	Created       []pr.Profile
}

func (r *ProfileRepo) Create(profile pr.Profile) error {
	if r.ErrOnCreate != nil {
		return r.ErrOnCreate
	}

	r.Created = append(r.Created, profile)
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
