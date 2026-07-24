package user

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	r "github.com/lejeunel/go-image-annotator/entities/role"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
	. "maragu.dev/gomponents"
)

var listUsersFields = []string{"id/email", "roles", "groups", "admin", "actions"}

type ListPresenter struct {
	b.PaginatedListBuilder
	Writer io.Writer
	e.ErrorPresenter
}

func NewListPresenter(w http.ResponseWriter, p b.PaginatedListBuilder) ListPresenter {
	return ListPresenter{p, w, e.NewErrorPresenter(w)}
}
func (p ListPresenter) SuccessListUsers(r list.Response) {
	p.SetPagination(r.Pagination, rt.AdminUsers)
	for _, user := range r.Users {
		row := MakeRow(user)
		p.AddRow(row)
	}

	p.PaginatedListBuilder.AddMarkdownPreamble(preamble)
	p.Render(p.Writer)
}

type ViewPresenter struct {
	io.Writer
	e.ErrorPresenter
}

func NewViewPresenter(w http.ResponseWriter) ViewPresenter {
	return ViewPresenter{Writer: w, ErrorPresenter: e.NewErrorPresenter(w)}
}
func (p ViewPresenter) SuccessFindUser(user u.User) {
	MakeRow(user).Render(p.Writer)
}

type DeletePresenter struct {
	io.Writer
	e.ErrorPresenter
}

func NewDeletePresenter(w http.ResponseWriter) DeletePresenter {
	return DeletePresenter{Writer: w, ErrorPresenter: e.NewErrorPresenter(w)}
}
func (p DeletePresenter) SuccessFindUser(user u.User) {
	b.RenderConfirmDeleteRow(len(listUsersFields),
		user.Id,
		"user",
		rt.AddQueryParams(User, "id", user.Id),
		rt.AddQueryParams(User, "id", user.Id, "mode", "view"),
		p.Writer)
}

type EditPresenter struct {
	io.Writer
	e.ErrorPresenter
	groups []g.Group
	roles  []r.Role
	user   u.User
}

func NewEditPresenter(w http.ResponseWriter) EditPresenter {
	return EditPresenter{Writer: w, ErrorPresenter: e.NewErrorPresenter(w)}
}

func (p *EditPresenter) SuccessFindUser(user u.User) {
	p.user = user
}
func (p *EditPresenter) SuccessListGroups(groups []g.Group) {
	p.groups = groups
}
func (p *EditPresenter) SuccessListRoles(roles []r.Role) {
	p.roles = roles
}
func (p EditPresenter) Render(w io.Writer) {
	endpoint := rt.AddQueryParams(User, "id", p.user.Id)
	form := bf.NewHTMXInlineFormBuilder(p.user.Id, len(listUsersFields), endpoint)
	form.AddTitle(fmt.Sprintf("Editing %v", p.user.Id))
	groupSelect := form.AddSelectableCombobox("Groups", "groups")
	for _, grp := range p.groups {
		groupSelect.AddField(grp.Name, slices.Contains(p.user.Groups, grp.Name))
	}
	roleSelect := form.AddSelectableCombobox("Roles", "roles")
	for _, role := range p.roles {
		roleSelect.AddField(role.Name, slices.Contains(p.user.Roles, role.Name))
	}
	form.Render(p.Writer)
}

func MakeRow(user u.User) tb.Row {
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(rt.AddQueryParams(User, "id", user.Id, "mode", "edit"))
	actions.SetConfirmDelete(rt.AddQueryParams(User, "id", user.Id,
		"mode", "confirm-delete"))
	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(user.Id)))
	row.AddCell(tb.NewCell(Text(strings.Join(user.Roles, ", "))))
	row.AddCell(tb.NewCell(Text(strings.Join(user.Groups, ", "))))
	row.AddCell(tb.NewCell(Text(strconv.FormatBool(user.IsAdmin))))
	row.AddCell(tb.NewCell(actions.Build()))
	return row

}
