package authorizers

import (
	"context"
)

type TestingIdentityProvider struct {
	Entitlements_ []string
	Groups_       []string
	Username_     string
	Email_        string
}

func (p *TestingIdentityProvider) Entitlements(ctx context.Context) ([]string, error) {
	return p.Entitlements_, nil
}

func (p *TestingIdentityProvider) Groups(ctx context.Context) ([]string, error) {
	return p.Groups_, nil
}

func (p *TestingIdentityProvider) Username(ctx context.Context) (string, error) {
	return p.Username_, nil
}

func (p *TestingIdentityProvider) Email(ctx context.Context) (string, error) {

	return p.Email_, nil
}
