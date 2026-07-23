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
		NewUserListPresenter(w, s.Page))
}

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	s.Users.Find.Execute(r.Context(), r.URL.Query().Get("id"), NewUserPresenter(w, r.URL.Query().Get("mode")))
}
func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXCreateFormBuilder(User, createUserTargetDiv)
	b.AddTitle("Create a new User")
	b.AddTextField("Email", "Email", "id", bf.WithRequired())
	b.AddCheckbox("Admin?", "Admin", "admin")
	b.Render(w)
}
