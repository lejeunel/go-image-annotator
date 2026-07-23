package user

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	rt "github.com/lejeunel/go-image-annotator/routes"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
	. "maragu.dev/gomponents"
)

var listUsersFields = []string{"id/email", "roles", "groups", "admin", "actions"}
var UserPage = "Users"

type ListUsersPresenter struct {
	b.PaginatedListBuilder
	Writer io.Writer
	e.WebPageErrorPresenter
}

func NewListUsersPresenter(w http.ResponseWriter, p b.PaginatedListBuilder) ListUsersPresenter {
	return ListUsersPresenter{p, w, e.NewErrorPresenter(w)}
}
func (p ListUsersPresenter) SuccessListUsers(r list.Response) {
	p.SetPagination(r.Pagination, rt.AdminUsers)
	for _, user := range r.Users {
		row := MakeListUserRow(user)
		p.AddRow(row)
	}
	p.Render(p.Writer)
}
func (p ListUsersPresenter) SuccessFindUser(user u.User) {
	MakeListUserRow(user).Render(p.Writer)
}

func (s *Server) ListUsers(w http.ResponseWriter, r *http.Request) {
	s.Page.SetUserIdentity(r.Context()).SetHTMLTitle("Users").SetTitle("Users")
	s.Page.ActivateSidebarEntry(UserPage)
	s.Page.AddCreationButton("Create", rt.CreateUserForm, createUserTargetDiv)
	s.Users.List.Execute(r.Context(), pa.PaginationParams{PageSize: s.DefaultPageSize, Page: pg.GetPageFromRequest(r)},
		NewListUsersPresenter(w, s.Page))
}

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	switch r.URL.Query().Get("mode") {
	case "edit":
		endpoint := rt.AddQueryParams(rt.User, "id", id)
		b := bf.NewHTMXInlineFormBuilder(id, len(listUsersFields), endpoint)
		b.AddTitle(fmt.Sprintf("Editing %v", id))
		b.Render(w)
	case "confirm-delete":
		b.RenderConfirmDeleteRow(len(listUsersFields),
			id,
			"user",
			rt.AddQueryParams(rt.User, "id", id),
			rt.AddQueryParams(rt.User, "id", id, "mode", "view"),
			w)
	default:
		s.Users.Find.Execute(r.Context(), id, NewListUsersPresenter(w, s.Page))
	}
}
func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXCreateFormBuilder(rt.User, createUserTargetDiv)
	b.AddTitle("Create a new User")
	b.AddTextField("Email", "Email", "id", bf.WithRequired())
	b.AddCheckbox("Admin?", "Admin", "admin")
	b.Render(w)
}
func MakeListUserRow(user u.User) tb.Row {
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(rt.AddQueryParams(rt.User, "id", user.Id, "mode", "edit"))
	actions.SetConfirmDelete(rt.AddQueryParams(rt.User, "id", user.Id,
		"mode", "confirm-delete"))

	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(user.Id)))
	row.AddCell(tb.NewCell(Text(strings.Join(user.Roles, ", "))))
	row.AddCell(tb.NewCell(Text(strings.Join(user.Groups, ", "))))
	row.AddCell(tb.NewCell(Text(strconv.FormatBool(user.IsAdmin))))
	row.AddCell(tb.NewCell(actions.Build()))
	return row

}
