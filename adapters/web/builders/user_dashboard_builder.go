package builders

import (
	"context"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"strings"
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
	User *u.User
}

func (b *UserDashboardBuilder) SetUserIdentityFromContext(ctx context.Context) *UserDashboardBuilder {
	id := u.IdentityFromContext(ctx)
	b.User = id
	return b
}

func MakeCard(node Node) Node {
	return Article(
		Class("group rounded-radius flex max-w-md flex-col border border-outline bg-surface-alt p-1 text-on-surface dark:border-outline-dark dark:bg-surface-dark-alt dark:text-on-surface-dark"),
		P(Class("text-pretty text-sm"), node),
	)
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
	apiTokenFrame, err := templatesFiles.ReadFile("templates/api_token_frame.html")
	if err != nil {
		return Text(err.Error())
	}
	apiTokenFrameStr := string(apiTokenFrame)

	APIToken := Div(Class("mt-2"), H3(Text("API Token")),
		P(Class("text-sm text-on-surface dark:text-on-surface-dark"),
			Text("Generate a secret token to authenticate your API requests. ")),
		Raw(apiTokenFrameStr))

	return Div(H2(Text("Profile")), MakeCard(profile), APIToken)
}

func NewUserDashboardBuilder() UserDashboardBuilder {
	return UserDashboardBuilder{}
}
