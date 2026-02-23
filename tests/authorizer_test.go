package tests

import (
	"context"
	a "datahub/app/authorizer"
	"testing"
)

func TestAuthorizerNoEmailBuildsPseudoEmail(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "username", "bob")
	auth := a.NewAuthorizer()
	email, err := auth.IdentityProvider.Email(ctx)
	want_email := "bob@no-email.com"
	AssertNoError(t, err)
	if email != want_email {
		t.Fatalf("expected to retrieve pseudo email. Wanted %v, got %v",
			want_email, email)
	}
}
