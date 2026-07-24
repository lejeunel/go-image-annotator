package auth

import (
	"fmt"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	rt "github.com/lejeunel/go-image-annotator/routes"
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
	if err := s.SessionManager.PasswordLogin(
		r.Context(), email, password); err != nil {
		htmx.NotifyError(w, "Login password", "wrong password")
		// w.Header().Set("HX-Trigger", string(pl))
		// w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	htmx.NotifySuccessPayloadAndRedirect(w, "Login", fmt.Sprintf("Successfully logged-in as %v", email), rt.HomePageUrl)
}
func (s Server) Login(w http.ResponseWriter, r *http.Request) {
	s.LoginPageBuilder.Render(w)
}
