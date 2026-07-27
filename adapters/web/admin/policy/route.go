package policy

import (
	_ "embed"
	"github.com/go-chi/chi/v5"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"io"
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

type ViewPresenter struct {
	b.PageBuilder
	io.Writer
	e.ErrorPresenter
}

func NewViewPresenter(w http.ResponseWriter, p b.PageBuilder) ViewPresenter {
	return ViewPresenter{p, w, e.NewErrorPresenter(w)}
}
func (p ViewPresenter) SuccessReadPolicy(policies string) {
	p.AddMarkdownPreamble(preamble)
	p.SetContent(
		Textarea(
			Class(`w-120 h-90 rounded-lg border-2 border-blue-500 p-3
         focus:outline-none focus:ring-2 focus:ring-blue-400
         resize-none`),
			Text(policies)))
	p.Render(p.Writer)
}

func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	s.Page.SetUserIdentity(r.Context())
	s.Itrs.Read.Execute(r.Context(), NewViewPresenter(w, s.Page))
}
