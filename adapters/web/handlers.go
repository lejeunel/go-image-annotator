package web

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/shared/html"
)

func RegisterWebPages(mux *http.ServeMux, server Server, b html.PageBuilder) {
	mux.Handle("/", HomePageHandlerFunc(b))
	mux.HandleFunc("/login/{provider}", server.HandleLogin)
	mux.HandleFunc("/logout", server.HandleLogout)
	mux.HandleFunc("/callback/{provider}", server.HandleAuthCallback)

	mux.HandleFunc("/user", server.User)

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
