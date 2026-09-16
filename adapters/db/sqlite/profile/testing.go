package profile

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

func CreateProfile(
	repo ProfileRepo,
	name string,
	group *string,
	labels *[]string,
) (*pr.Profile, error) {
	p := pr.NewProfile(pr.NewProfileId(), name,
		pr.WithDescription("a-description"))

	if group != nil {
		p.Group = group
	}
	if labels != nil {
		p.Labels = *labels
	}
	if err := repo.Create(p); err != nil {
		return nil, err
	}
	return &p, nil
}
