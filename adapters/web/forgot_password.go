package web

import (
	"log/slog"
	. "maragu.dev/gomponents"
	"net/http"

	fpw "github.com/lejeunel/go-image-annotator/use-cases/user/forgot-password"
)

type PasswordResetNotifier interface {
	Notify(r fpw.Response)
}

type VoidPasswordResetNotifier struct {
	slog.Logger
}

func (n VoidPasswordResetNotifier) Notify(r fpw.Response) {
	n.Logger.Info("notifying password reset token", "notification", r)
}

type PasswordResetPresenter struct {
	PasswordResetNotifier
	slog.Logger
	w http.ResponseWriter
}

func (p PasswordResetPresenter) Success(r fpw.Response) {
	p.Notify(r)
	p.redirect()
}
func (p PasswordResetPresenter) Error(err error) {
	p.Logger.Error(err.Error())
	p.redirect()
}
func (p PasswordResetPresenter) redirect() {
	Text("If the provided email exists in our database, you will receive an email shortly.").Render(p.w)
}

func EmitPasswordResetTokenHandlerFunc(itr fpw.Interactor, logger slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := PasswordResetPresenter{VoidPasswordResetNotifier{logger}, logger, w}
		r.ParseForm()
		itr.Execute(r.Context(), r.FormValue("email"), p)
	}
}
