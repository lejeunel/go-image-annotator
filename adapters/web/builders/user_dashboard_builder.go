package builders

import (
	_ "embed"
	"fmt"
	"strings"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	rt "github.com/lejeunel/go-image-annotator/routes"
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

type UserDashboardBuilder struct {
	PageBuilder
}

func (b *UserDashboardBuilder) Build() *UserDashboardBuilder {
	if b.User == nil {
		b.SetError(fmt.Errorf("failed build user dashboard: user identity has not been set"))
		return b
	}
	rows := []UserInfoRow{{Name: "Email", Value: b.User.Id}}
	if b.User.IsAdmin {
		rows = append(rows, UserInfoRow{Name: "Is admin", Value: "yes"})
	}
	rows = append(rows, UserInfoRow{Name: "Groups", Value: strings.Join(b.User.Groups, ", ")})
	rows = append(rows, UserInfoRow{Name: "Roles", Value: strings.Join(b.User.Roles, ", ")})
	profile := Table(Class("text-left text-sm text-on-surface dark:text-on-surface-dark"),
		Map(rows, func(r UserInfoRow) Node {
			return r.Render()
		}),
	)
	APIToken := Div(Class("mt-2"), H3(Text("API token")),
		P(Class("text-sm text-on-surface dark:text-on-surface-dark"),
			Text("Generate a secret token to authenticate your API requests. ")),
		Raw(cmp.ApiTokenFrame))

	changePassword := Div(Class("mt-2"), H3(Text("Reset password")),
		cmp.MakeCard(Form(
			Attr(fmt.Sprintf(`hx-post=%v`, rt.ChangePassword)),
			Class("bg-surface-alt/50 dark:bg-surface-dark-alt/50 p-4 rounded-lg shadow-md mb-4"),
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

	content := Div(cmp.MakeCard(profile), APIToken, changePassword)
	b.SetContent(content)
	return b
}

func NewUserDashboardBuilder(b PageBuilder) UserDashboardBuilder {
	return UserDashboardBuilder{PageBuilder: b}
}
