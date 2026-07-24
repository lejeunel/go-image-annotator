package user

import (
	_ "embed"
	"net/http"

	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

//go:embed preamble.md
var preamble string

var UserPage = "Users"

func (s *Server) ListUsers(w http.ResponseWriter, r *http.Request) {
	s.Page.SetUserIdentity(r.Context()).SetHTMLTitle("Users").SetTitle("Users")
	s.Page.ActivateSidebarEntry(UserPage)
	s.Page.AddCreationButton("Create", CreateUserForm, createUserTargetDiv)
	s.Users.List.Execute(r.Context(), pa.PaginationParams{PageSize: s.DefaultPageSize, Page: pg.GetPageFromRequest(r)},
		NewListPresenter(w, s.Page))
}

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	switch r.URL.Query().Get("mode") {
	case "edit":
		p := NewEditPresenter(w)
		s.Roles.List.Execute(r.Context(), &p)
		s.Groups.List.Execute(r.Context(), &p)
		s.Users.Find.Execute(r.Context(), id, &p)
		p.Render(w)
	case "confirm-delete":
		s.Users.Find.Execute(r.Context(), id, NewDeletePresenter(w))
	default:
		s.Users.Find.Execute(r.Context(), id, NewViewPresenter(w))
	}
}
func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXCreateFormBuilder(User, createUserTargetDiv)
	b.AddTitle("Create a new User")
	b.AddTextField("email", "Email", "email", bf.WithRequired())
	b.AddCheckbox("is_admin", "Admin", "admin")
	b.Render(w)
}
