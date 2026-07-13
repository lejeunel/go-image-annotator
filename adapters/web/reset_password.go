package web

import (
	"log/slog"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	rt "github.com/lejeunel/go-image-annotator/routes"
	rfpw "github.com/lejeunel/go-image-annotator/use-cases/user/reset-forgotten-password"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type PasswordResetPresenter struct {
	slog.Logger
	w http.ResponseWriter
	HTMXErrorPresenter
}

func (p PasswordResetPresenter) Success() {
	Div(P(Text("Password changed successfully!")),
		P(Text("Proceed to "), cmp.MakeTextLink(rt.Home, "login"), Text("."))).Render(p.w)
}
func (p PasswordResetPresenter) Error(err error) {
	p.HTMXErrorPresenter.Error(err)
}

func ResetPassword(itr rfpw.Interactor, logger slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		token := r.URL.Query().Get("token")
		pw1 := r.FormValue("password")
		pw2 := r.FormValue("password-repeat")
		p := PasswordResetPresenter{logger, w, NewHTMXErrorPresenter("Resetting password", w)}
		r.ParseForm()
		itr.Execute(r.Context(), rfpw.Request{Token: token, FirstPassword: pw1, SecondPassword: pw2}, p)
	}
}

func ResetPasswordForm(builder b.ResetPasswordBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		builder.SetToken(r.URL.Query().Get("token"))
		builder.Render(w)
	}

}
