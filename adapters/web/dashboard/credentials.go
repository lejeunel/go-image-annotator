package dashboard

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	cpw "github.com/lejeunel/go-image-annotator/use-cases/user/change-password"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

//go:embed credentials-preamble.md
var credentialsPreamble string

func makeSectionTitle(title string) Node {
	return Div(Class("text-lg font-bold"), Text(title))
}

func RenderCredentialsPage(ctx context.Context, pb b.PageBuilder, w io.Writer) {
	APIToken := Div(Class("mt-2"), makeSectionTitle("API token"),
		P(Class("text-sm text-on-surface dark:text-on-surface-dark"),
			Text("Generate a secret token to authenticate your API requests. ")),
		Raw(cmp.ApiTokenFrame))

	form := bf.NewHTMXFormBuilder(ChangePasswordUrl, "change-password-form")
	form.AddTextField(
		CurrentPasswordFieldName,
		"Current password",
		bf.WithRequired(),
		bf.WithHidden(),
	)
	form.AddTextField(NewPasswordFieldName, "New password", bf.WithRequired(), bf.WithHidden())
	form.AddTextField(
		RepeatPasswordFieldName,
		"New password (bis)",
		bf.WithRequired(),
		bf.WithHidden(),
	)
	changePassword := Div(Class("mt-2"), makeSectionTitle("Reset password"),
		form.Build())

	content := Div(
		Class("flex flex-col w-120"),
		Div(cmp.Separator, APIToken, cmp.Separator, changePassword),
	)
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
	fmt.Println(r.Form)
	s.ChangePasswordItr.Execute(r.Context(), cpw.Request{
		Id:              user.Id,
		CurrentPassword: r.FormValue(CurrentPasswordFieldName),
		FirstPassword: r.FormValue(
			NewPasswordFieldName,
		),
		SecondPassword: r.FormValue(RepeatPasswordFieldName),
	},
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
