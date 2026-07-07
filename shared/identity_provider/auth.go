package provider

import (
	"net/http"

	s "github.com/lejeunel/go-image-annotator/shared/session"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type AuthHandler interface {
	OAuthLogin(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
	OAuthCallback(http.ResponseWriter, *http.Request)
	PasswordLogin(http.ResponseWriter, *http.Request)
}

type MyAuthHandler struct {
	s.SessionManager
}

func (p MyAuthHandler) OAuthLogin(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}
func (p MyAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	p.SessionManager.Logout(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (p MyAuthHandler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		if err := p.SessionManager.FinishOAuthLogin(r.Context(), user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
func (p MyAuthHandler) PasswordLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if err := p.SessionManager.PasswordLogin(r.Context(), email, password); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type OAuthProviderConfig struct {
	Key         string
	Secret      string
	CallbackURL string
}

func SetupForGoogle(cfg OAuthProviderConfig) {
	goth.UseProviders(google.New(cfg.Key, cfg.Secret, cfg.CallbackURL))
}

func NewAuthHandler(m s.SessionManager) MyAuthHandler {
	return MyAuthHandler{m}
}
