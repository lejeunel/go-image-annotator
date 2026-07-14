package auth

import (
	"net/http"
)

func (s Server) ForgotPasswordForm(w http.ResponseWriter, r *http.Request) {
	s.ForgotPasswordBuilder.Render(w)
}
