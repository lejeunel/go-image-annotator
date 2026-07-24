package collection

import (
	_ "embed"
	"net/http"

	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

//go:embed preamble.md
var preamble string

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	s.FindItr.Execute(r.Context(),
		r.URL.Query().Get("name"),
		NewCollectionPresenter(w, r.URL.Query().Get("mode")))
}
func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXCreateFormBuilder(Collection, createCollectionTargetDiv)
	b.AddTitle("Create a new collection")
	b.AddTextField("name", "Name", "name", bf.WithRequired())
	b.AddTextField("description", "Description", "description")
	b.Render(w)
}
func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.ListItr.Execute(r.Context(), pa.PaginationParams{PageSize: s.DefaultPageSize, Page: pg.GetPageFromRequest(r)},
		NewListPresenter(w, s.PageBuilder))
}
