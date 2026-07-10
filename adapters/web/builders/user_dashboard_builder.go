package builders

import (
	"context"
	_ "embed"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"strings"
)

//go:embed components/api_token_frame.html
var ApiTokenFrame string

type UserInfoRow struct {
	Name  string
	Value string
}

func (r UserInfoRow) Render() Node {
	return Tr(Td(Class("py-2 px-2 font-bold"), Text(r.Name)),
		Td(Class("py-2 px-2"), Text(r.Value)))
}

type UserDashboardBuilder struct {
	User *u.User
}

func (b *UserDashboardBuilder) SetUserIdentityFromContext(ctx context.Context) *UserDashboardBuilder {
	id := u.IdentityFromContext(ctx)
	b.User = id
	return b
}

func (b *UserDashboardBuilder) Build() Node {
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
		Raw(ApiTokenFrame))

	return Div(cmp.MakeCard(profile), APIToken)
}

func NewUserDashboardBuilder() UserDashboardBuilder {
	return UserDashboardBuilder{}
}
