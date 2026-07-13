package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
)

func LoginPageHandlerFunc(builder b.LoginPageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		builder.Render(w)
	}
}

func HomePageHandlerFunc(pb b.PageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pb.SetUserIdentity(r.Context())
		web.MakeHomePage(pb, w)
	}
}
func APIDocsHandlerFunc(specsPath string, pb b.PageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		web.APIDocsPage(r.Context(), specsPath, *pb.SetActive(cmp.APIDocsPageActive), w)
	}
}
