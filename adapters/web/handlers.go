package web

import (
	"net/http"
)

func RegisterWebPages(mux *http.ServeMux, server Server) {
	mux.HandleFunc("/", HomePageHandler)
	mux.HandleFunc("/collections", server.ListCollections)
	mux.HandleFunc("/images", server.ListImages)
	mux.HandleFunc("/labels", server.ListLabels)
	mux.HandleFunc("/image", server.ViewImage)
	mux.HandleFunc("/ui/submit-box", server.SubmitBox)
}
