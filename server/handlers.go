package server

import (
	_ "embed"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func APIDocsHandlerFunc(specsPath string, pb b.PageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content := Div(Raw("<redoc spec-url='/api/openapi.yaml'></redoc>"),
			Script(Src("/static/redoc.standalone.js")))
		pb.SetUserIdentity(r.Context())
		pb.SetExpanded()
		pb.SetHTMLTitle("API Docs")
		pb.SetActiveSection(cmp.APIDocsPageActive)
		pb.SetContent(Div(Class("bg-white"), content))
		pb.Render(w)
	}
}
