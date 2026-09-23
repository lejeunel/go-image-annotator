package web

import (
	"io"
	"net/http"

	_ "embed"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
)

//go:embed text.md
var text string

func HandlerFunc(pb b.PageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pb.SetUserIdentity(r.Context())
		pb.SetHTMLTitle("Home")
		MakeHomePage(pb, w)
	}
}

func MakeHomePage(pb b.PageBuilder, w io.Writer) {
	pb.SetTitle("Home")
	pb.SetActiveSection(cmp.HomePageActive)
	pb.AddMarkdownPreamble(text)
	pb.Render(w)
}
