package profile

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

func CreateProfile(repo ProfileRepo, name string) (*pr.Profile, error) {
	p := pr.NewProfile(pr.NewProfileId(), name,
		pr.WithDescription("a-description"))
	if err := repo.Create(p); err != nil {
		return nil, err
	}
	return &p, nil
}
