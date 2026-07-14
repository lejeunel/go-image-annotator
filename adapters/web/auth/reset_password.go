package auth

import (
	"log/slog"
	"net/http"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	rt "github.com/lejeunel/go-image-annotator/routes"
	rfpw "github.com/lejeunel/go-image-annotator/use-cases/user/reset-forgotten-password"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type PasswordResetPresenter struct {
	slog.Logger
	w http.ResponseWriter
	htmx.ErrorPresenter
}

func (p PasswordResetPresenter) Success() {
	Div(P(Text("Password changed successfully!")),
		P(Text("Proceed to "), cmp.MakeTextLink(rt.Home, "login"), Text("."))).Render(p.w)
}
func (p PasswordResetPresenter) Error(err error) {
	p.ErrorPresenter.Error(err)
}

func (s Server) ResetPassword(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token := r.URL.Query().Get("token")
	pw1 := r.FormValue("password")
	pw2 := r.FormValue("password-repeat")
	pres := PasswordResetPresenter{s.Logger, w, htmx.NewErrorPresenter("Resetting password", w)}
	r.ParseForm()
	s.ResetPasswordItr.Execute(r.Context(), rfpw.Request{Token: token, FirstPassword: pw1, SecondPassword: pw2}, pres)
}

func (s Server) ResetPasswordForm(w http.ResponseWriter, r *http.Request) {
	s.ResetPasswordBuilder.SetToken(r.URL.Query().Get("token"))
	s.ResetPasswordBuilder.Render(w)

}
