package web

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	"net/http"
)

func RegisterWebPages(mux *http.ServeMux, server Server, b b.PageBuilder) {
	mux.Handle("/", HomePageHandlerFunc(b))
	mux.HandleFunc("/user-dashboard", server.UserDashboard)
	mux.HandleFunc("/ui/new-api-token", server.NewAPIToken)

	mux.HandleFunc("/collections", server.ListCollections)
	mux.HandleFunc("/images", server.ListImages)
	mux.HandleFunc("/labels", server.ListLabels)
	mux.HandleFunc("/image", server.ViewImage)

	mux.HandleFunc("/ui/annotate/submit-box", server.SubmitBox)
	mux.HandleFunc("/ui/annotate/update-box", server.UpdateBox)
	mux.HandleFunc("/ui/annotate/submit-polygon", server.SubmitPolygon)
	mux.HandleFunc("/ui/annotate/update-polygon", server.UpdatePolygon)
	mux.HandleFunc("/ui/annotate/submit-label", server.SubmitLabel)
	mux.HandleFunc("/ui/annotate/annotation-panel", server.MakeAnnotationPanel)
	mux.HandleFunc("/ui/annotate/annotations", server.GetRegionsAsJSON)
	mux.HandleFunc("/ui/annotate/remove-annotation", server.DeleteAnnotation)
	mux.HandleFunc("/ui/annotate/set-label", server.SetLabel)
}

func MakeAuthMux(h ip.AuthHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/login/password", h.PasswordLogin)
	mux.HandleFunc("/login/{provider}", h.OAuthLogin)
	mux.HandleFunc("/callback/{provider}", h.OAuthCallback)
	mux.HandleFunc("/logout", h.Logout)
	return mux
}
