package authorizers

import (
	"context"
	e "datahub/errors"
	"fmt"
	"net/http"
	"slices"
)

func AuthentikHeadersMiddleware(next http.Handler) http.Handler {
	// We wrap our anonymous function, and cast it to a http.HandlerFunc
	// Because our function signature matches ServeHTTP(w, r), this allows
	// our function (type) to implicitly satisify the http.Handler interface.
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			username := r.Header.Get("X-Authentik-Username")
			ctx := context.WithValue(r.Context(), "username", username)

			email := r.Header.Get("X-Authentik-Email")
			ctx2 := context.WithValue(ctx, "email", email)

			entitlements := r.Header.Get("X-Authentik-Entitlements")
			ctx3 := context.WithValue(ctx2, "entitlements", entitlements)

			groups := r.Header.Get("X-Authentik-Groups")
			ctx4 := context.WithValue(ctx3, "groups", groups)

			newReq := r.WithContext(ctx4)

			next.ServeHTTP(w, newReq)
		})
}

type Authorizer struct {
	IdentityProvider IdentityProvider
}

func NewAuthorizer() *Authorizer {
	return &Authorizer{
		IdentityProvider: &HeaderIdentityProvider{EntitlementsSeparator: "|"}}
}

func (a *Authorizer) wantGroupMembership(ctx context.Context, targetGroup string) error {
	currentUserGroups, err := a.IdentityProvider.Groups(ctx)
	if err != nil {
		return fmt.Errorf("fetching groups: %w", err)
	}
	if !slices.Contains(currentUserGroups, targetGroup) {
		return fmt.Errorf("authorizing user with group memberships %v to contribute a resource needing group membership %v: %w",
			currentUserGroups, targetGroup, e.ErrGroupOwnership)
	}
	return nil

}
func (a *Authorizer) isAdmin(ctx context.Context) (bool, error) {
	entitlements, err := a.IdentityProvider.Entitlements(ctx)
	if err != nil {
		return false, fmt.Errorf("checking for admin entitlement given entitlements %v: %w",
			entitlements, err)
	}
	if slices.Contains(entitlements, "admin") {
		return true, nil
	}
	return false, nil

}

func (a *Authorizer) wantEntitlement(ctx context.Context, neededEntitlement string) error {
	entitlements, err := a.IdentityProvider.Entitlements(ctx)
	if err != nil {
		return fmt.Errorf("authorizing user with entitlements %v but need entitlements %v: %w",
			entitlements, neededEntitlement, err)
	}

	if !slices.Contains(entitlements, neededEntitlement) {
		return fmt.Errorf("authorizing user with entitlements %v, needing %v: %w",
			entitlements, neededEntitlement, e.ErrEntitlement)
	}
	return nil

}

func (a *Authorizer) WantToUpdateCollection(ctx context.Context, targetGroup string) error {
	isAdmin, _ := a.isAdmin(ctx)
	if isAdmin {
		return nil
	}

	return a.WantToContributeImages(ctx, targetGroup)
}

func (a *Authorizer) WantToContributeImages(ctx context.Context, targetGroup string) error {
	isAdmin, _ := a.isAdmin(ctx)
	if isAdmin {
		return nil
	}
	if err := a.wantEntitlement(ctx, "im-contrib"); err != nil {
		return err
	}
	if err := a.wantGroupMembership(ctx, targetGroup); err != nil {
		return err
	}

	return nil

}
func (a *Authorizer) WantToContributeLabels(ctx context.Context) error {
	isAdmin, _ := a.isAdmin(ctx)
	if isAdmin {
		return nil
	}
	if err := a.wantEntitlement(ctx, "annotation-contrib"); err != nil {
		return err
	}
	return nil

}

func (a *Authorizer) WantToContributeAnnotations(ctx context.Context, imageGroup string) error {
	isAdmin, _ := a.isAdmin(ctx)
	if isAdmin {
		return nil
	}
	if err := a.wantEntitlement(ctx, "annotation-contrib"); err != nil {
		return err
	}
	if err := a.wantGroupMembership(ctx, imageGroup); err != nil {
		return err
	}

	return nil
}

func (a *Authorizer) WantToDeleteCollectionOrItsContent(ctx context.Context, collectionGroup string) error {
	isAdmin, _ := a.isAdmin(ctx)
	if isAdmin {
		return nil
	}
	if err := a.wantEntitlement(ctx, "im-contrib"); err != nil {
		return err
	}
	if err := a.wantGroupMembership(ctx, collectionGroup); err != nil {
		return err
	}

	return nil
}

func (a *Authorizer) WantToContributeLocation(ctx context.Context, group string) error {
	isAdmin, _ := a.isAdmin(ctx)
	if isAdmin {
		return nil
	}
	if err := a.wantEntitlement(ctx, "im-contrib"); err != nil {
		return err
	}
	if err := a.wantGroupMembership(ctx, group); err != nil {
		return err
	}

	return nil
}
