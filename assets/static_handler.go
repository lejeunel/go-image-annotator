package assets

import (
	"net/http"
)

func RegisterStaticFiles(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("assets/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
}
