package auth

import (
	"log/slog"
	"net/http"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	rt "github.com/lejeunel/go-image-annotator/routes"
	fpw "github.com/lejeunel/go-image-annotator/use-cases/user/forgot-password"
)

type NotifyPasswordResetPresenter struct {
	Notifier
	slog.Logger
	w       http.ResponseWriter
	baseURL string
}

func (p NotifyPasswordResetPresenter) Success(r fpw.Response) {
	url := rt.AddQueryParams(p.baseURL, "token", r.PasswordResetToken)
	p.Notify(Notification{Email: r.Email, URL: url.String()})
	p.redirect()
}

func (p NotifyPasswordResetPresenter) Error(err error) {
	p.Logger.Error(err.Error())
	p.redirect()
}

func (p NotifyPasswordResetPresenter) redirect() {
	Div(
		Text("If the provided email exists in our database, you will receive an email shortly."),
	).Render(p.w)
}

func (s Server) NotifyPasswordReset(w http.ResponseWriter, r *http.Request) {
	p := NotifyPasswordResetPresenter{
		s.Notifier, s.Logger, w,
		s.baseURL + rt.ResetPasswordFormUrl,
	}
	r.ParseForm()
	s.requestTokenItr.Execute(r.Context(), r.FormValue("email"), p)
}
