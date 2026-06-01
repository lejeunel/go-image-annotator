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
	mux.HandleFunc("/ui/submit-label", server.SubmitLabel)
	mux.HandleFunc("/ui/annotation-panel", server.MakeHTMLAnnotationPanel)
	mux.HandleFunc("/ui/annotations", server.GetAnnotationsAsJSON)
	mux.HandleFunc("/ui/remove-annotation", server.DeleteAnnotation)
	mux.HandleFunc("/ui/update-box", server.UpdateBox)
	mux.HandleFunc("/ui/set-label", server.SetLabel)
}
