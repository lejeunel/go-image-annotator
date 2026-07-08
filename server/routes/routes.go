package routes

import (
	"fmt"
	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	"github.com/lejeunel/go-image-annotator/adapters/web"
	as "github.com/lejeunel/go-image-annotator/assets"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"

	"github.com/go-chi/chi/v5"
	"net/http"
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

func RouteAuth(r chi.Router, h ip.AuthHandler,
	loginPage http.HandlerFunc, forgotPasswordPage http.HandlerFunc,
	sessionMiddleware func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(sessionMiddleware)
		r.HandleFunc(LoginWithPassword, h.PasswordLogin)
		r.HandleFunc(LoginOAuth, h.OAuthLogin)
		r.HandleFunc(CallbackOAuth, h.OAuthCallback)
		r.HandleFunc(Logout, h.Logout)
	})

	r.Handle(Login, loginPage)
	r.Handle(ForgotPassword, forgotPasswordPage)
}

func RouteAPIDocs(r chi.Router, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	mwChain := chi.Chain(mws...)
	r.Method(http.MethodGet, APIDocs, mwChain.HandlerFunc(h))
}

func RouteWebPages(r chi.Router, s web.Server, home http.HandlerFunc,
	mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(Home, home)
		r.Get(UserDashboard, s.UserDashboard)
		r.Get(NewAPIToken, s.NewAPIToken)

		r.Get(Collections, s.ListCollections)
		r.Get(Images, s.ListImages)
		r.Get(Labels, s.ListLabels)
		r.Get(Image, s.ViewImage)

		r.Post(SubmitBox, s.SubmitBox)
		r.Put(UpdateBox, s.UpdateBox)
		r.Post(SubmitPolygon, s.SubmitPolygon)
		r.Put(UpdatePolygon, s.UpdatePolygon)
		r.Post(SubmitLabel, s.SubmitLabel)
		r.Get(AnnotationPanel, s.MakeAnnotationPanel)
		r.Get(Annotations, s.GetRegionsAsJSON)
		r.Delete(RemoveAnnotation, s.DeleteAnnotation)
		r.Post(SetLabel, s.SetLabel)
	})
}
func RouteAPI(r chi.Router, apiServer api.Server, mws ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Handle(APIRoot, api.Handler(&apiServer))
	})
}

func RouteAPISpecs(r chi.Router) {
	r.HandleFunc(APISpecs,
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			w.Write(as.Openapiyaml)
		})
}

func RouteStaticFiles(r chi.Router) {
	fs := http.FileServer(http.Dir("assets/static"))
	r.Handle(fmt.Sprintf("%v/*", StaticRoot), http.StripPrefix(StaticRoot, fs))
}
