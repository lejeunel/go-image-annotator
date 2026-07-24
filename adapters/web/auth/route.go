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
		r.HandleFunc(rt.LoginOAuthUrl, s.OAuthLogin)
		r.HandleFunc(rt.CallbackOAuthUrl, s.OAuthCallback)
		r.HandleFunc(rt.LogoutUrl, s.Logout)
		r.Post(rt.LoginPageUrl, s.PasswordLogin)
	})

	r.Get(rt.LoginPageUrl, s.Login)
	r.Get(rt.ForgotPasswordFormUrl, s.ForgotPasswordForm)
	r.Post(rt.NotifyPasswordResetUrl, s.NotifyPasswordReset)
	r.Get(rt.ResetPasswordFormUrl, s.ResetPasswordForm)
	r.Post(rt.ResetPasswordUrl, s.ResetPassword)
}
