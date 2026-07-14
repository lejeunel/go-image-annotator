package auth

import (
	rt "github.com/lejeunel/go-image-annotator/routes"

	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Server) Route(r chi.Router,
	sessionMiddleware func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(sessionMiddleware)
		r.HandleFunc(rt.LoginWithPassword, s.PasswordLogin)
		r.HandleFunc(rt.LoginOAuth, s.OAuthLogin)
		r.HandleFunc(rt.CallbackOAuth, s.OAuthCallback)
		r.HandleFunc(rt.Logout, s.Logout)
	})

	r.Get(rt.Login, s.Login)
	r.Get(rt.ForgotPasswordForm, s.ForgotPasswordForm)
	r.Post(rt.NotifyPasswordReset, s.NotifyPasswordReset)
	r.Get(rt.ResetPasswordForm, s.ResetPasswordForm)
	r.Post(rt.ResetPassword, s.ResetPassword)
}
