package cli

import (
	"github.com/lejeunel/go-image-annotator/entities/principal"
)

type PrincipalProvider struct{}

func (p *PrincipalProvider) Provide() (*principal.Principal, error) {
	return &principal.Principal{}, nil
}
