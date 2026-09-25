package user

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

func (s *Server) ListUsers(w http.ResponseWriter, r *http.Request) {
	s.PaginatedListBuilder.SetUserIdentity(r.Context()).SetHTMLTitle("Users").SetTitle("Users")
	s.PaginatedListBuilder.ActivateSidebarEntry(PageName)
	s.PaginatedListBuilder.AddCreationButton("Create", CreateUserFormUrl, createUserTargetDiv)
	s.Users.List.Execute(
		r.Context(),
		pa.PaginationParams{PageSize: s.DefaultPageSize, Page: pg.GetPageFromRequest(r)},
		NewListPresenter(w, s.PaginatedListBuilder, s.RowUrl),
	)
}

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(resourceUrlFieldName)
	s.RowUrl.SetId(id)
	switch r.URL.Query().Get("mode") {
	case b.ModeEdit.String():
		p := NewEditPresenter(w, s.PageBuilder, s.RowUrl)
		s.Roles.List.Execute(r.Context(), &p)
		s.Groups.List.Execute(r.Context(), &p)
		s.Users.Find.Execute(r.Context(), id, &p)
		p.Render(w)
	case b.ModeConfirmDelete.String():
		s.Users.Find.Execute(r.Context(), id, NewDeletePresenter(w, s.PageBuilder, s.RowUrl))
	default:
		s.Users.Find.Execute(r.Context(), id, NewViewPresenter(w, s.PageBuilder, s.RowUrl))
	}
}

func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXCreateFormBuilder(UserUrl, createUserTargetDiv)
	b.AddTitle("Create a new User")
	b.AddTextField(createEmailFieldName, "Email", bf.WithRequired())
	b.AddCheckbox(createIsAdminFieldName, "Admin")
	b.Render(w)
}
