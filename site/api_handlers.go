package site

import (
	_ "embed"
	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	"net/http"
)

//go:embed openapi.yaml
var openapiyaml []byte

func RegisterAPIDocs(mux *http.ServeMux, specsPath, docsPath string) {
	mux.HandleFunc(docsPath,
		func(w http.ResponseWriter, r *http.Request) {
			if err := APIDocsPage(r.Context(), specsPath, docsPath).Render(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
}

func RegisterAPISpecs(mux *http.ServeMux, path string) {
	mux.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			w.Write(openapiyaml)
		})
}

func RegisterAPIEndpoints(mux *http.ServeMux, server api.Server, path string) {
	api.HandlerFromMuxWithBaseURL(&server, mux, path)
}
