package auth

import (
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func (s Server) OAuthLogin(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (s Server) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		if err := s.SessionManager.FinishOAuthLogin(r.Context(), user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

type OAuthProviderConfig struct {
	Key         string
	Secret      string
	CallbackURL string
}

func SetupForGoogle(cfg OAuthProviderConfig) {
	goth.UseProviders(google.New(cfg.Key, cfg.Secret, cfg.CallbackURL))
}
