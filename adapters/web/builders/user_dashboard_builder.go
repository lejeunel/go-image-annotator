package builders

import (
	_ "embed"
	"fmt"
	"strings"

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

	APIToken := Div(Class("mt-2"), H3(Text("API Token")),
		P(Class("text-sm text-on-surface dark:text-on-surface-dark"),
			Text("Generate a secret token to authenticate your API requests. ")),
		Raw(cmp.ApiTokenFrame))

	content := Div(cmp.MakeCard(profile), APIToken)
	b.SetContent(content)
	return b
}

func NewUserDashboardBuilder(b PageBuilder) UserDashboardBuilder {
	return UserDashboardBuilder{PageBuilder: b}
}
