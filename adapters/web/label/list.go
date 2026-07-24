package label

import (
	_ "embed"
	"net/http"

	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

//go:embed preamble.md
var preamble string

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	s.FindItr.Execute(r.Context(),
		r.URL.Query().Get("name"),
		NewLabelPresenter(w, r.URL.Query().Get("mode")))
}
func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.ListItr.Execute(r.Context(),
		pag.PaginationParams{PageSize: s.DefaultPageSize, Page: pg.GetPageFromRequest(r)},
		NewListPresenter(w, s.PageBuilder))
}
