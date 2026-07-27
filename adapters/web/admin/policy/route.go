package policy

import (
	_ "embed"
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
)

//go:embed preamble.md
var preamble string

func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.AdminPoliciesUrl, s.Edit)
	})
}

func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	s.Page.SetUserIdentity(r.Context())
	// TODO find current policy using interactor and append to textarea

	s.Page.AddMarkdownPreamble(preamble)
	s.Page.SetContent(Textarea(Text("hello")))
	s.Page.Render(w)
}
