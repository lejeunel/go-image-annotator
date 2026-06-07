package site

import (
	_ "embed"
	"fmt"
	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	"net/http"
)

//go:embed openapi.yaml
var openapiyaml []byte

func RegisterAPIDocs(mux *http.ServeMux, server api.Server, apiPath string) {
	api.HandlerFromMuxWithBaseURL(&server, mux, fmt.Sprintf("/%v", apiPath))
	RegisterAPISpecs(mux, apiPath)
}

func RegisterAPISpecs(mux *http.ServeMux, apiPath string) {
	specsPath := fmt.Sprintf("/%v/openapi.yml", apiPath)
	docsPath := fmt.Sprintf("/%v/docs", apiPath)
	mux.HandleFunc(specsPath,
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			w.Write(openapiyaml)
		})
	mux.Handle(docsPath, APIDocsHandler(specsPath, apiPath))
}

func APIDocsHandler(specURL, apiPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := APIDocsPage(specURL, apiPath).Render(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
