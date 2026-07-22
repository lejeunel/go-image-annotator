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
var LoginWithPassword = "/auth/login/password"
var LoginOAuth = "/auth/login/{provider}"
var CallbackOAuth = "/auth/callback/{provider}"
var ForgotPasswordForm = "/auth/forgot-password"
var NotifyPasswordReset = "/auth/notify-password-reset"
var ResetPasswordForm = "/auth/reset-password-form"
var ResetPassword = "/auth/reset-password"
var Logout = "/auth/logout"

var UserDashboard = "/user-dashboard"
var NewAPIToken = "/ui/new-api-token"

var Home = "/"
var Collections = "/collections"
var Images = "/images"
var Labels = "/labels"

var AnnotateImage = "/ui/annotate/image"
var SubmitBox = "/ui/annotate/submit-box"
var UpdateBox = "/ui/annotate/update-box"
var SubmitPolygon = "/ui/annotate/submit-polygon"
var UpdatePolygon = "/ui/annotate/update-polygon"
var SubmitImageLabel = "/ui/annotate/submit-label"
var AnnotationPanel = "/ui/annotate/annotation-panel"
var Annotations = "/ui/annotate/annotations"
var RemoveAnnotation = "/ui/annotate/remove-annotation"
var SetLabel = "/ui/annotate/set-label"

var Image = "/ui/image"

var Collection = "/ui/collection"
var CreateCollectionForm = "/ui/collection/new"

var Label = "/ui/label"
var CreateLabelForm = "/ui/label/new"

var Admin = "/admin"
var AdminUsers = "/admin/users"
var AdminGroups = "/admin/groups"
var AdminRoles = "/admin/roles"
var AdminPolicies = "/admin/policies"

var User = "/ui/user"
var CreateUserForm = "/ui/user/new"

func MakeOAuthCallbackURL(baseURL string, provider string) string {
	return baseURL + strings.ReplaceAll(CallbackOAuth, "{provider}", provider)
}
func MakeOAuthLoginURL(provider string) string {
	return strings.ReplaceAll(LoginOAuth, "{provider}", provider)
}

func MakeAnnotateImageURL(imageId, collection string) string {
	return fmt.Sprintf("%v?id=%v&collection=%v", AnnotateImage, imageId, collection)
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
