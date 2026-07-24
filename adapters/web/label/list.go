package label

import (
	_ "embed"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

//go:embed preamble.md
var preamble string

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get(resourceUrlFieldName)
	s.RowURL.SetId(name)
	switch r.URL.Query().Get("mode") {
	case b.ModeEdit.String():
		s.FindItr.Execute(r.Context(), name, NewEditPresenter(w, s.RowURL))
	case b.ModeConfirmDelete.String():
		s.FindItr.Execute(r.Context(), name, NewDeletePresenter(w, s.RowURL))
	default:
		s.FindItr.Execute(r.Context(), name, NewViewPresenter(w, s.RowURL))
	}

}
func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.ListItr.Execute(r.Context(),
		pag.PaginationParams{PageSize: s.DefaultPageSize, Page: pg.GetPageFromRequest(r)},
		NewListPresenter(w, s.PageBuilder, s.RowURL))
}
