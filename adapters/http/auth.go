package http

import (
	"context"
	"github.com/lejeunel/go-image-annotator/entities/principal"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type PrincipalProvider struct {
	c   context.Context
	key string
}

func NewPrincipalProvider(c context.Context) PrincipalProvider {
	return PrincipalProvider{c: c, key: "principal"}
}

func (p PrincipalProvider) Provide() (*principal.Principal, error) {
	pr, ok := p.c.Value(p.key).(principal.Principal)
	if !ok {
		return nil, e.ErrPrincipalProvider
	}
	return &pr, nil
}
