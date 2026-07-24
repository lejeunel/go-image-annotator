package server

import (
	"fmt"
	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	as "github.com/lejeunel/go-image-annotator/assets"
	rt "github.com/lejeunel/go-image-annotator/routes"

	"github.com/go-chi/chi/v5"
	"net/http"
)

func RouteAPIDocs(r chi.Router, h http.HandlerFunc, mws ...func(http.Handler) http.Handler) {
	mwChain := chi.Chain(mws...)
	r.Method(http.MethodGet, rt.APIDocsUrl, mwChain.HandlerFunc(h))
}

func RouteWebPages(r chi.Router, home http.HandlerFunc,
	mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.HomePageUrl, home)
	})
}
func RouteAPI(r chi.Router, apiServer api.Server, mws ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(mws...)
		handler := api.HandlerWithOptions(&apiServer, api.StdHTTPServerOptions{
			BaseURL: rt.APIRootUrl,
		})
		r.Mount(rt.APIRootUrl, handler)
	})
}

func RouteAPISpecs(r chi.Router) {
	r.HandleFunc(rt.APISpecsUrl,
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			w.Write(as.Openapiyaml)
		})
}

func RouteStaticFiles(r chi.Router) {
	fs := http.FileServer(http.Dir("assets/static"))
	r.Handle(fmt.Sprintf("%v/*", rt.StaticRootUrl), http.StripPrefix(rt.StaticRootUrl, fs))
}
