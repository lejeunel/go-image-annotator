package site

import (
	_ "embed"
	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	"net/http"
)

//go:embed openapi.yaml
var openapiyaml []byte

func RegisterAPI(mux *http.ServeMux, server api.Server, docsPath string, specsPath string) {
	api.HandlerFromMuxWithBaseURL(&server, mux, "/api")
	RegisterAPISpecs(mux, docsPath, specsPath)
}

func RegisterAPISpecs(mux *http.ServeMux, docsPath, specsPath string) {
	mux.HandleFunc(specsPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(openapiyaml)
	})
	mux.Handle(docsPath, APIDocsHandler(specsPath))
}

func APIDocsHandler(specURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := APIDocsPage(specURL).Render(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
