package web

import (
	"net/http"
	"strings"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type UserInfoRow struct {
	Name  string
	Value string
}

func (r UserInfoRow) Render() Node {
	return Tr(Td(Class("py-2 px-2 font-bold"), Text(r.Name)),
		Td(Class("py-2 px-2"), Text(r.Value)))

}

func (s *Server) User(w http.ResponseWriter, r *http.Request) {
	p := s.PageBuilder
	p.SetUserIdentityFromContext(r.Context())
	user := p.User
	p.SetTitle("User Dashboard")
	rows := []UserInfoRow{{Name: "Email", Value: user.Id}}
	if user.IsAdmin {
		rows = append(rows, UserInfoRow{Name: "Is admin", Value: "yes"})
	}
	rows = append(rows, UserInfoRow{Name: "Groups", Value: strings.Join(user.Groups, ", ")})
	rows = append(rows, UserInfoRow{Name: "Roles", Value: strings.Join(user.Roles, ", ")})
	info := Table(Class("text-left text-sm text-on-surface dark:text-on-surface-dark"),
		Map(rows, func(r UserInfoRow) Node {
			return r.Render()
		}),
	)

	p.SetContent(info)
	p.Render(w)

}
