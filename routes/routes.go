package routes

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	APIRootUrl  = "/api"
	APISpecsUrl = "/api/openapi.yaml"
	APIDocsUrl  = "/api/docs"

	StaticRootUrl = "/static"

	LoginPageUrl           = "/auth/login"
	LoginOAuthUrl          = "/auth/login/{provider}"
	CallbackOAuthUrl       = "/auth/callback/{provider}"
	ForgotPasswordFormUrl  = "/auth/forgot-password"
	NotifyPasswordResetUrl = "/auth/notify-password-reset"
	ResetPasswordFormUrl   = "/auth/reset-password-form"
	ResetPasswordUrl       = "/auth/reset-password"
	LogoutUrl              = "/auth/logout"

	HomePageUrl      = "/"
	CollectionsUrl   = "/collections"
	ImagesUrl        = "/images"
	LabelsUrl        = "/labels"
	UserDashboardUrl = "/user-dashboard"

	AdminUrl         = "/admin"
	AdminUsersUrl    = "/admin/users"
	AdminGroupsUrl   = "/admin/groups"
	AdminRolesUrl    = "/admin/roles"
	AdminPoliciesUrl = "/admin/policies"
)

func MakeOAuthCallbackURL(baseURL string, provider string) string {
	return baseURL + strings.ReplaceAll(CallbackOAuthUrl, "{provider}", provider)
}
func MakeOAuthLoginURL(provider string) string {
	return strings.ReplaceAll(LoginOAuthUrl, "{provider}", provider)
}

func MakeAnnotateImageURL(baseURL, imageId, collection string) string {
	return fmt.Sprintf("%v?id=%v&collection=%v", baseURL, imageId, collection)
}

func MakeImagesURL(collection string) string {
	return fmt.Sprintf("%v?collection=%v", ImagesUrl, collection)
}

func AddQueryParams(baseURL string, params ...string) url.URL {
	u, _ := url.Parse(baseURL)
	q := u.Query()
	for i := 0; i+1 < len(params); i += 2 {
		q.Add(params[i], params[i+1])
	}

	u.RawQuery = q.Encode()
	return *u
}
