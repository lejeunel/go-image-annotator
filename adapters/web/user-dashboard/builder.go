package user

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"strings"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

//go:embed preamble.md
var preamble string

type UserInfoRow struct {
	Name  string
	Value string
}

func (r UserInfoRow) Render() Node {
	return Tr(Td(Class("py-2 px-2 font-bold"), Text(r.Name)),
		Td(Class("py-2 px-2"), Text(r.Value)))
}

func makeSectionTitle(title string) Node {
	return Div(Class("text-lg font-bold"), Text(title))

}
func RenderDashboard(ctx context.Context, pb b.PageBuilder, w io.Writer) {
	if pb.User == nil {
		pb.SetError(fmt.Errorf("failed build user dashboard: user identity has not been set"))
		pb.Render(w)
		return
	}
	rows := []UserInfoRow{{Name: "Email", Value: pb.User.Id}}
	if pb.User.IsAdmin {
		rows = append(rows, UserInfoRow{Name: "Is admin", Value: "yes"})
	}
	rows = append(rows, UserInfoRow{Name: "Groups", Value: strings.Join(pb.User.Groups, ", ")})
	rows = append(rows, UserInfoRow{Name: "Roles", Value: strings.Join(pb.User.Roles, ", ")})
	profile := Table(Class("text-left text-sm text-on-surface dark:text-on-surface-dark"),
		Map(rows, func(r UserInfoRow) Node {
			return r.Render()
		}),
	)
	APIToken := Div(Class("mt-2"), makeSectionTitle("API token"),
		P(Class("text-sm text-on-surface dark:text-on-surface-dark"),
			Text("Generate a secret token to authenticate your API requests. ")),
		Raw(cmp.ApiTokenFrame))

	changePassword := Div(Class("mt-2"), makeSectionTitle("Reset password"),
		cmp.MakeCard(Form(
			Attr(fmt.Sprintf(`hx-post=%v`, ChangePassword)),
			Class("m-2"),
			Label(For("Current password"), Text("Current password"), Class(st.FormLabel)),
			Input(Type("password"), ID("password-current"), Name("password-current"), Required(), Class(st.FormInput)),
			Label(For("New password"), Text("New password"), Class(st.FormLabel)),
			Input(Type("password"), ID("password"), Name("password"), Required(), Class(st.FormInput)),
			Label(For("New password (repeat)"), Text("New password (repeat)"), Class(st.FormLabel)),
			Input(Type("password"), ID("password-repeat"), Name("password-repeat"), Required(), Class(st.FormInput)),
			Button(Type("submit"),
				Text("Submit"),
				Class(st.SuccessButton)),
		),
		))

	separator := Div(Class("h-px bg-outline dark:bg-outline-dark"))
	content := Div(Class("flex flex-col w-120"), Div(cmp.MakeCard(profile), separator, APIToken, separator, changePassword))
	pb.SetUserIdentity(ctx)
	pb.SetActiveSection(cmp.NoPageActive)
	pb.SetTitle("User Dashboard")
	pb.SetHTMLTitle("Dashboard")
	pb.AddMarkdownPreamble(preamble)
	pb.SetContent(content)
	pb.Render(w)
}
