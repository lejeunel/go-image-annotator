package authorizers

import (
	"context"
	e "datahub/errors"
	"fmt"
	"strings"
)

type IdentityProvider interface {
	Entitlements(context.Context) ([]string, error)
	Groups(context.Context) ([]string, error)
	Username(context.Context) (string, error)
	Email(context.Context) (string, error)
}

type HeaderIdentityProvider struct {
	EntitlementsSeparator string
}

func (p *HeaderIdentityProvider) Entitlements(ctx context.Context) ([]string, error) {
	entitlements := ctx.Value("entitlements")
	if entitlements == nil {
		return []string{}, fmt.Errorf("found no entitlements field in context: %w", e.ErrIdentity)
	}

	return strings.Split(entitlements.(string), p.EntitlementsSeparator), nil

}

func (p *HeaderIdentityProvider) Username(ctx context.Context) (string, error) {
	username := ctx.Value("username")
	if username == nil {
		return "", fmt.Errorf("getting username: %w", e.ErrIdentity)
	}
	return username.(string), nil
}

func (p *HeaderIdentityProvider) Email(ctx context.Context) (string, error) {
	email := ctx.Value("email")
	if email == nil {
		username, err := p.Username(ctx)
		if err != nil {
			return "", fmt.Errorf("attempted to build pseudo email from username: %w", e.ErrIdentity)
		}
		return fmt.Sprintf("%v@no-email.com", username), nil
	}
	return email.(string), nil
}

func (p *HeaderIdentityProvider) Groups(ctx context.Context) ([]string, error) {
	groups := ctx.Value("groups")
	if groups == nil {
		return []string{}, fmt.Errorf("got no group assignements: %w", e.ErrIdentity)
	}
	return strings.Split(groups.(string), p.EntitlementsSeparator), nil
}
