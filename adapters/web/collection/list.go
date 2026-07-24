package collection

import (
	_ "embed"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

//go:embed preamble.md
var preamble string

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get(resourceUrlFieldName)
	s.RowURL.SetId(name)
	switch r.URL.Query().Get("mode") {
	case b.ModeEdit.String():
		s.FindItr.Execute(r.Context(), name,
			NewEditPresenter(w, s.RowURL))
	case b.ModeConfirmDelete.String():
		s.FindItr.Execute(r.Context(), name,
			NewDeletePresenter(w, s.RowURL))
	default:
		s.FindItr.Execute(r.Context(), name,
			NewViewPresenter(w, s.RowURL))
	}
}
func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXCreateFormBuilder(CollectionUrl, createCollectionTargetDiv)
	b.AddTitle("Create a new collection")
	b.AddTextField(createNameFieldName, "Name", bf.WithRequired())
	b.AddTextField(createDescriptionFieldName, "Description")
	b.Render(w)
}
func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.ListItr.Execute(r.Context(), pa.PaginationParams{PageSize: s.DefaultPageSize, Page: pg.GetPageFromRequest(r)},
		NewListPresenter(w, s.PageBuilder, s.RowURL))
}
