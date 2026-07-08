package server

import (
	"github.com/lejeunel/go-image-annotator/adapters/web"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"net/http"
)

func LoginPageHandlerFunc(builder b.LoginPageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		builder.Render(w)
	}
}
func ForgotPasswordHandlerFunc(builder b.ForgotPasswordBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		builder.Render(w)
	}
}

func HomePageHandlerFunc(pb b.PageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pb.SetUserIdentityFromContext(r.Context())
		web.MakeHomePage(pb).Render(w)
	}
}
func APIDocsHandlerFunc(specsPath string, pb b.PageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := web.APIDocsPage(r.Context(), specsPath, *pb.SetActive(b.APIDocsPageActive)).Render(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
