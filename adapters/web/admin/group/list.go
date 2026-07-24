package group

import (
	_ "embed"
	"net/http"

	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
)

//go:embed preamble.md
var preamble string

var GroupPage = "Groups"

func (s *Server) ListGroups(w http.ResponseWriter, r *http.Request) {
	s.Page.SetUserIdentity(r.Context()).SetHTMLTitle("Groups").SetTitle("Groups")
	s.Page.ActivateSidebarEntry(GroupPage)
	s.Page.AddCreationButton("Create", CreateUserForm, createGroupTargetDiv)
	s.Groups.List.Execute(r.Context(), NewListPresenter(w, s.Page))
}

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	switch r.URL.Query().Get("mode") {
	case "edit":
		s.Groups.Find.Execute(r.Context(), name, NewEditPresenter(w))
	case "confirm-delete":
		s.Groups.Find.Execute(r.Context(), name, NewDeletePresenter(w))
	default:
		s.Groups.Find.Execute(r.Context(), name, NewViewPresenter(w))
	}
}
func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXCreateFormBuilder(GroupRow, createGroupTargetDiv)
	b.AddTitle("Create a new group")
	b.AddTextField("name", "Name", "name", bf.WithRequired())
	b.AddTextField("description", "Description", "description")
	b.Render(w)
}
