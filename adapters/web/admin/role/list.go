package role

import (
	_ "embed"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
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
		s.Roles.Find.Execute(r.Context(), name, NewEditPresenter(w, s.RowUrl))
	case b.ModeConfirmDelete.String():
		s.Roles.Find.Execute(r.Context(), name, NewDeletePresenter(w, s.RowUrl))
	default:
		s.Roles.Find.Execute(r.Context(), name, NewViewPresenter(w, s.RowUrl))
	}
}
func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXCreateFormBuilder(RoleRowUrl, createRoleTargetDiv)
	b.AddTitle("Create a new role")
	b.AddTextField(createNameFieldName, "Name", bf.WithRequired())
	b.AddTextField(createDescriptionFieldName, "Description")
	b.Render(w)
}
