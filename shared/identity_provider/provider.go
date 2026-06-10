package provider

import (
	"context"
	"net/http"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	s "github.com/lejeunel/go-image-annotator/shared/session"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func IdentityFromContext(ctx context.Context) *u.User {
	v := ctx.Value(u.UserContextKey)
	if v == nil {
		return nil
	}
	user, ok := v.(*u.User)
	if !ok {
		return nil
	}

	return user
}

type OAuthHandler interface {
	HandleLogin(http.ResponseWriter, *http.Request)
	HandleLogout(http.ResponseWriter, *http.Request)
	HandleAuthCallback(http.ResponseWriter, *http.Request)
}

type GothIdentityHandler struct {
	s.SessionManager
}

func (p GothIdentityHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (p GothIdentityHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	p.SessionManager.Logout(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (p GothIdentityHandler) HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		if err := p.SessionManager.Login(r.Context(), user.Email); err != nil {
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

func NewGothIdentityHandler(m s.SessionManager) GothIdentityHandler {
	return GothIdentityHandler{m}
}
