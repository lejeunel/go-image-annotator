package server

import (
	"fmt"
	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	"github.com/lejeunel/go-image-annotator/adapters/web"
	as "github.com/lejeunel/go-image-annotator/assets"
	rt "github.com/lejeunel/go-image-annotator/routes"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"

	"github.com/go-chi/chi/v5"
	"net/http"
)

func RouteAuth(r chi.Router, h ip.AuthHandler,
	loginPage http.HandlerFunc, forgotPasswordPage http.HandlerFunc,
	sessionMiddleware func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(sessionMiddleware)
		r.HandleFunc(rt.LoginWithPassword, h.PasswordLogin)
		r.HandleFunc(rt.LoginOAuth, h.OAuthLogin)
		r.HandleFunc(rt.CallbackOAuth, h.OAuthCallback)
		r.HandleFunc(rt.Logout, h.Logout)
	})

	r.Handle(rt.Login, loginPage)
	r.Handle(rt.ForgotPassword, forgotPasswordPage)
}

func RouteAPIDocs(r chi.Router, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	mwChain := chi.Chain(mws...)
	r.Method(http.MethodGet, rt.APIDocs, mwChain.HandlerFunc(h))
}

func RouteWebPages(r chi.Router, s web.Server, home http.HandlerFunc,
	mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.Home, home)
		r.Get(rt.UserDashboard, s.UserDashboard)
		r.Get(rt.NewAPIToken, s.NewAPIToken)

		r.Get(rt.Collections, s.ListCollections)
		r.Get(rt.Images, s.ListImages)
		r.Get(rt.Labels, s.ListLabels)
		r.Get(rt.Image, s.ViewImage)

		r.Post(rt.SubmitBox, s.SubmitBox)
		r.Put(rt.UpdateBox, s.UpdateBox)
		r.Post(rt.SubmitPolygon, s.SubmitPolygon)
		r.Put(rt.UpdatePolygon, s.UpdatePolygon)
		r.Post(rt.SubmitImageLabel, s.SubmitLabel)
		r.Get(rt.AnnotationPanel, s.MakeAnnotationPanel)
		r.Get(rt.Annotations, s.GetRegionsAsJSON)
		r.Delete(rt.RemoveAnnotation, s.DeleteAnnotation)
		r.Post(rt.SetLabel, s.SetLabel)

		r.Get(rt.Collection, s.GetCollection)
		r.Get(rt.CreateCollectionForm, s.CreateCollectionForm)
		r.Get(rt.ConfirmDeleteCollection, s.ConfirmDeleteCollection)
		r.Get(rt.EditCollectionForm, s.EditCollectionForm)
		r.Post(rt.Collection, s.CreateCollection)
		r.Delete(rt.Collection, s.DeleteCollection)
		r.Put(rt.Collection, s.EditCollection)

		r.Get(rt.Label, s.GetLabel)
		r.Get(rt.CreateLabelForm, s.CreateLabelForm)
		r.Get(rt.EditLabelForm, s.EditLabelForm)
		r.Put(rt.Label, s.EditLabel)
		r.Get(rt.ConfirmDeleteLabel, s.ConfirmDeleteLabel)
		r.Delete(rt.Label, s.DeleteLabel)
		r.Post(rt.Label, s.CreateLabel)
	})
}
func RouteAPI(r chi.Router, apiServer api.Server, mws ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(mws...)
		handler := api.HandlerWithOptions(&apiServer, api.StdHTTPServerOptions{
			BaseURL: rt.APIRoot,
		})
		r.Mount(rt.APIRoot, handler)
	})
}

func RouteAPISpecs(r chi.Router) {
	r.HandleFunc(rt.APISpecs,
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			w.Write(as.Openapiyaml)
		})
}

func RouteStaticFiles(r chi.Router) {
	fs := http.FileServer(http.Dir("assets/static"))
	r.Handle(fmt.Sprintf("%v/*", rt.StaticRoot), http.StripPrefix(rt.StaticRoot, fs))
}
