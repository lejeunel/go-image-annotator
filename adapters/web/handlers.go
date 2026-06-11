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

	mux.HandleFunc("/user-dashboard", server.UserDashboard)
	mux.HandleFunc("/ui/new-api-token", server.NewAPIToken)

	mux.HandleFunc("/collections", server.ListCollections)
	mux.HandleFunc("/images", server.ListImages)
	mux.HandleFunc("/labels", server.ListLabels)
	mux.HandleFunc("/image", server.ViewImage)

	mux.HandleFunc("/ui/annotate/submit-box", server.SubmitBox)
	mux.HandleFunc("/ui/annotate/submit-label", server.SubmitLabel)
	mux.HandleFunc("/ui/annotate/annotation-panel", server.MakeHTMLAnnotationPanel)
	mux.HandleFunc("/ui/annotate/annotations", server.GetAnnotationsAsJSON)
	mux.HandleFunc("/ui/annotate/remove-annotation", server.DeleteAnnotation)
	mux.HandleFunc("/ui/annotate/update-box", server.UpdateBox)
	mux.HandleFunc("/ui/annotate/set-label", server.SetLabel)
}
