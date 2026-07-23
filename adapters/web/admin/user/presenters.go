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
	u "github.com/lejeunel/go-image-annotator/entities/user"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
	. "maragu.dev/gomponents"
)

var listUsersFields = []string{"id/email", "roles", "groups", "admin", "actions"}

type UserListPresenter struct {
	b.PaginatedListBuilder
	Writer io.Writer
	e.WebPageErrorPresenter
}

func NewUserListPresenter(w http.ResponseWriter, p b.PaginatedListBuilder) UserListPresenter {
	return UserListPresenter{p, w, e.NewErrorPresenter(w)}
}
func (p UserListPresenter) SuccessListUsers(r list.Response) {
	p.SetPagination(r.Pagination, rt.AdminUsers)
	for _, user := range r.Users {
		row := MakeListUserRow(user)
		p.AddRow(row)
	}

	p.PaginatedListBuilder.AddMarkdownPreamble(preamble)
	p.Render(p.Writer)
}

type UserPresenter struct {
	io.Writer
	e.WebPageErrorPresenter
	successFindUser func(u.User)
}

func NewUserPresenter(w http.ResponseWriter, mode string) UserPresenter {
	p := UserPresenter{Writer: w, WebPageErrorPresenter: e.NewErrorPresenter(w)}
	switch mode {
	case "edit":
		p.successFindUser = p.renderEditUserForm
	case "confirm-delete":
		p.successFindUser = p.renderConfirmDelete
	default:
		p.successFindUser = p.renderView
	}
	return p
}
func (p UserPresenter) SuccessFindUser(user u.User) {
	p.successFindUser(user)
}
func (p UserPresenter) renderView(user u.User) {
	MakeListUserRow(user).Render(p.Writer)
}
func (p UserPresenter) renderEditUserForm(user u.User) {
	endpoint := rt.AddQueryParams(User, "id", user.Id)
	form := bf.NewHTMXInlineFormBuilder(user.Id, len(listUsersFields), endpoint)
	form.AddTitle(fmt.Sprintf("Editing %v", user.Id))
	form.Render(p.Writer)
}
func (p UserPresenter) renderConfirmDelete(user u.User) {
	b.RenderConfirmDeleteRow(len(listUsersFields),
		user.Id,
		"user",
		rt.AddQueryParams(User, "id", user.Id),
		rt.AddQueryParams(User, "id", user.Id, "mode", "view"),
		p.Writer)
}

func MakeListUserRow(user u.User) tb.Row {
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
