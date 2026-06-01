package auth

import (
	p "github.com/lejeunel/go-image-annotator/entities/principal"
)

type PrincipalProvider interface {
	Provide() (*p.Principal, error)
}

type PassThroughAuth struct {
}

func (a PassThroughAuth) CreateCollection(p PrincipalProvider, group string) error {
	return nil
}
