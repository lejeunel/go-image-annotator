package role

import (
	_ "embed"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
)

//go:embed preamble.md
var preamble string

func (s *Server) ListRoles(w http.ResponseWriter, r *http.Request) {
	s.Page.SetUserIdentity(r.Context()).SetHTMLTitle("Roles").SetTitle("Roles")
	s.Page.AddCreationButton("Create", CreateRoleForm, createRoleTargetDiv)
	s.Roles.List.Execute(r.Context(), NewListPresenter(w, s.Page, s.RowUrl))
}

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get(resourceUrlFieldName)
	s.RowUrl.SetId(name)
	switch r.URL.Query().Get("mode") {
	case b.ModeEdit.String():
		p := NewEditPresenter(w, s.RowUrl)
		p.SetValidMethods(s.Auth.ListMethods())
		s.Roles.Find.Execute(r.Context(), name, &p)
		p.Render()
	case b.ModeConfirmDelete.String():
		s.Roles.Find.Execute(r.Context(), name, NewDeletePresenter(w, s.Page.PageBuilder, s.RowUrl))
	default:
		p := NewViewPresenter(w, s.Page.PageBuilder, s.RowUrl)
		s.Roles.Find.Execute(r.Context(), name, &p)
	}
}
