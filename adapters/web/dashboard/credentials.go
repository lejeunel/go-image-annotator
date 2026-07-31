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
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	cpw "github.com/lejeunel/go-image-annotator/use-cases/user/change-password"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

//go:embed credentials-preamble.md
var credentialsPreamble string

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
func RenderCredentialsPage(ctx context.Context, pb b.PageBuilder, w io.Writer) {
	if pb.User == nil {
		pb.SetError(fmt.Errorf("failed build user dashboard: user identity has not been set"))
		pb.Render(w)
		return
	}
	rows := []UserInfoRow{{Name: "Email", Value: pb.User.Id}}
	if pb.User.IsAdmin() {
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
			Attr(fmt.Sprintf(`hx-post=%v`, ChangePasswordUrl)),
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

	content := Div(Class("flex flex-col w-120"), Div(cmp.MakeCard(profile), cmp.Separator, APIToken, cmp.Separator, changePassword))
	pb.SetActiveSection(cmp.NoPageActive)
	pb.AddMarkdownPreamble(credentialsPreamble)
	pb.SetContent(content)
	pb.Render(w)
}

func (s *Server) Credentials(w http.ResponseWriter, r *http.Request) {
	s.SetUserIdentity(r.Context())
	s.SetTitle(CredentialsPageName)
	s.ActivateSidebarEntry(CredentialsPageName)
	s.SetHTMLTitle(CredentialsPageName)
	RenderCredentialsPage(r.Context(), s.PageBuilder, w)
}
func (s *Server) NewAPIToken(w http.ResponseWriter, r *http.Request) {
	user := u.IdentityFromContext(r.Context())
	if user == nil {
		http.Error(w, "failed getting user identity", http.StatusForbidden)
	}
	s.RenewAPITokenItr.Execute(r.Context(),
		user.Id, cmp.NewAPITokenPresenter(w))
}
func (s *Server) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := u.IdentityFromContext(r.Context())
	if user == nil {
		http.Error(w, "failed getting user identity", http.StatusForbidden)
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	s.ChangePasswordItr.Execute(r.Context(), cpw.Request{Id: user.Id, CurrentPassword: r.FormValue("password-current"),
		FirstPassword: r.FormValue("password"), SecondPassword: r.FormValue("password-repeat")},
		NewChangePasswordPresenter(w))

}

type ChangePasswordPresenter struct {
	writer http.ResponseWriter
	task   string
	htmx.ErrorPresenter
}

func NewChangePasswordPresenter(w http.ResponseWriter) ChangePasswordPresenter {
	task := "Change password"
	return ChangePasswordPresenter{w, task, htmx.NewErrorPresenter(task, w)}
}
func (p ChangePasswordPresenter) Success() {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, "Successfully changed password")
}
