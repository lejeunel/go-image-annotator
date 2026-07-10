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
var ForgotPassword = "/auth/forgot-password"
var Logout = "/auth/logout"

var UserDashboard = "/user-dashboard"
var NewAPIToken = "/ui/new-api-token"

var Home = "/"
var Collections = "/collections"
var Images = "/images"
var Labels = "/labels"
var Image = "/image"

var SubmitBox = "/ui/annotate/submit-box"
var UpdateBox = "/ui/annotate/update-box"
var SubmitPolygon = "/ui/annotate/submit-polygon"
var UpdatePolygon = "/ui/annotate/update-polygon"
var SubmitImageLabel = "/ui/annotate/submit-label"
var AnnotationPanel = "/ui/annotate/annotation-panel"
var Annotations = "/ui/annotate/annotations"
var RemoveAnnotation = "/ui/annotate/remove-annotation"
var SetLabel = "/ui/annotate/set-label"

var Collection = "/ui/collection"
var ConfirmDeleteCollection = "/ui/collection/confirm-delete"
var EditCollectionForm = "/ui/collection/edit"
var CreateCollectionForm = "/ui/collection/new"

func MakeOAuthCallbackURL(baseURL string, provider string) string {
	return baseURL + strings.ReplaceAll(CallbackOAuth, "{provider}", provider)
}
func MakeOAuthLoginURL(provider string) string {
	return strings.ReplaceAll(LoginOAuth, "{provider}", provider)
}

func MakeImageURL(imageId, collection string) string {
	return fmt.Sprintf("%v?id=%v&collection=%v", Image, imageId, collection)
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
