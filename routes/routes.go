package routes

import (
	"fmt"
	"strings"
)

var APIRoot = "/api/"
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
var SubmitLabel = "/ui/annotate/submit-label"
var AnnotationPanel = "/ui/annotate/annotation-panel"
var Annotations = "/ui/annotate/annotations"
var RemoveAnnotation = "/ui/annotate/remove-annotation"
var SetLabel = "/ui/annotate/set-label"

func MakeOAuthCallbackURL(baseURL string, port int, provider string) string {
	return fmt.Sprintf("%v:%v%v/%v", baseURL, port, CallbackOAuth, provider)
}
func MakeOAuthLoginURL(provider string) string {
	return strings.ReplaceAll(LoginOAuth, "{provider}", provider)
}
