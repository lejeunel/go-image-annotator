package server

import (
	_ "embed"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func HomePageHandlerFunc(pb b.PageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pb.SetUserIdentity(r.Context())
		pb.SetHTMLTitle("Home")
		web.MakeHomePage(pb, w)
	}
}

func APIDocsHandlerFunc(specsPath string, pb b.PageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content := Div(Raw("<redoc spec-url='/api/openapi.yaml'></redoc>"),
			Script(Src("/static/redoc.standalone.js")))
		pb.SetUserIdentity(r.Context())
		pb.SetHTMLTitle("API Docs")
		pb.SetActiveSection(cmp.APIDocsPageActive)
		pb.SetContent(Div(Class("bg-white"), content))
		pb.Render(w)
	}
}
