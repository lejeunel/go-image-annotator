package server

import (
	"fmt"
	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	"github.com/lejeunel/go-image-annotator/adapters/web"
	as "github.com/lejeunel/go-image-annotator/assets"
	rt "github.com/lejeunel/go-image-annotator/routes"

	"github.com/go-chi/chi/v5"
	"net/http"
)

func RouteAPIDocs(r chi.Router, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	mwChain := chi.Chain(mws...)
	r.Method(http.MethodGet, rt.APIDocs, mwChain.HandlerFunc(h))
}

func RouteWebPages(r chi.Router, s web.Server, home http.HandlerFunc,
	mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.Home, home)

		r.Get(rt.AnnotateImage, s.AnnotateImage)

		r.Post(rt.SubmitBox, s.SubmitBox)
		r.Put(rt.UpdateBox, s.UpdateBox)
		r.Post(rt.SubmitPolygon, s.SubmitPolygon)
		r.Put(rt.UpdatePolygon, s.UpdatePolygon)
		r.Post(rt.SubmitImageLabel, s.SubmitLabel)
		r.Get(rt.AnnotationPanel, s.MakeAnnotationPanel)
		r.Get(rt.Annotations, s.GetRegionsAsJSON)
		r.Delete(rt.RemoveAnnotation, s.DeleteAnnotation)
		r.Post(rt.SetLabel, s.SetLabel)
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
