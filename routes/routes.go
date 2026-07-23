package routes

import (
	"fmt"
	"net/url"
	"strings"
)

var APIRoot = "/api"
var APISpecs = "/api/openapi.yaml"
var APIDocs = "/api/docs"

var StaticRoot = "/static"

var Login = "/auth/login"
var LoginOAuth = "/auth/login/{provider}"
var CallbackOAuth = "/auth/callback/{provider}"
var ForgotPasswordForm = "/auth/forgot-password"
var NotifyPasswordReset = "/auth/notify-password-reset"
var ResetPasswordForm = "/auth/reset-password-form"
var ResetPassword = "/auth/reset-password"
var Logout = "/auth/logout"

var Home = "/"
var Collections = "/collections"
var Images = "/images"
var Labels = "/labels"
var UserDashboard = "/user-dashboard"

var Image = "/ui/image"

var Admin = "/admin"
var AdminUsers = "/admin/users"
var AdminGroups = "/admin/groups"
var AdminRoles = "/admin/roles"
var AdminPolicies = "/admin/policies"

func MakeOAuthCallbackURL(baseURL string, provider string) string {
	return baseURL + strings.ReplaceAll(CallbackOAuth, "{provider}", provider)
}
func MakeOAuthLoginURL(provider string) string {
	return strings.ReplaceAll(LoginOAuth, "{provider}", provider)
}

func MakeAnnotateImageURL(baseURL, imageId, collection string) string {
	return fmt.Sprintf("%v?id=%v&collection=%v", baseURL, imageId, collection)
}

func MakeImagesURL(collection string) string {
	return fmt.Sprintf("%v?collection=%v", Images, collection)
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
