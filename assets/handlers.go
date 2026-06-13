package assets

import (
	_ "embed"
	"net/http"
)

//go:embed openapi.yaml
var openapiyaml []byte

func RegisterAPISpecs(mux *http.ServeMux, path string) {
	mux.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			w.Write(openapiyaml)
		})
}

func RegisterStaticFiles(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("assets/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
}
