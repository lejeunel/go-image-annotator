package dashboard

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"strings"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
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
func RenderProfilePage(ctx context.Context, pb b.PageBuilder, w io.Writer) {
	if pb.User == nil {
		pb.SetError(fmt.Errorf("failed build user dashboard: user identity has not been set"))
		pb.Render(w)
		return
	}
	rows := []UserInfoRow{{Name: "Email", Value: pb.User.Id}}
	rows = append(rows, UserInfoRow{Name: "Groups", Value: strings.Join(pb.User.Groups, ", ")})
	rows = append(rows, UserInfoRow{Name: "Roles", Value: strings.Join(pb.User.Roles, ", ")})
	profile := Table(Class("text-left text-sm text-on-surface dark:text-on-surface-dark"),
		Map(rows, func(r UserInfoRow) Node {
			return r.Render()
		}),
	)
	content := Div(Class("flex flex-col w-120"), Div(cmp.MakeCard(profile)))
	pb.SetActiveSection(cmp.NoPageActive)
	pb.SetContent(content)
	pb.Render(w)
}

func (s *Server) Profile(w http.ResponseWriter, r *http.Request) {
	s.SetUserIdentity(r.Context())
	s.SetTitle(ProfilePageName)
	s.ActivateSidebarEntry(ProfilePageName)
	s.SetHTMLTitle(ProfilePageName)
	RenderProfilePage(r.Context(), s.PageBuilder, w)
}
