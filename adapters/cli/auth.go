package cli

import (
	i "github.com/lejeunel/go-image-annotator/entities/identity"
)

type IdentityProvider struct{}

func (p *IdentityProvider) Provide() (*i.Identity, error) {
	return &i.Identity{}, nil
}
