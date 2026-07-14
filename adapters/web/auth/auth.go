package auth

import (
	"net/http"

	"github.com/markbates/goth/gothic"
)

func (s Server) Logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	s.SessionManager.Logout(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (s Server) PasswordLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if err := s.SessionManager.PasswordLogin(r.Context(), email, password); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (s Server) Login(w http.ResponseWriter, r *http.Request) {
	s.LoginPageBuilder.Render(w)
}
