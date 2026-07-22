package auth

import (
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

	if err := s.SessionManager.PasswordLogin(
		r.Context(),
		r.FormValue("email"),
		r.FormValue("password")); err != nil {
		pl, _ := htmx.NotifyError("Login password", "wrong password")
		w.Header().Set("HX-Trigger", string(pl))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.Header().Set("HX-Redirect", rt.Home)
}
func (s Server) Login(w http.ResponseWriter, r *http.Request) {
	s.LoginPageBuilder.Render(w)
}
