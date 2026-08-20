package auth

import (
	"log/slog"
	"net/http"
	"os"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	rt "github.com/lejeunel/go-image-annotator/routes"
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

func MaybeSetupGoogle(pb *b.LoginPageBuilder, baseURL string, logger slog.Logger) {
	id := os.Getenv("GOIA_GOOGLE_CLIENT_ID")
	secret := os.Getenv("GOIA_GOOGLE_CLIENT_SECRET")
	if (id != "") && (secret != "") {
		logger.Info("setting up google auth")
		pb.AddOAuthProvider(ProviderNameGoogle, rt.MakeOAuthLoginURL(ProviderNameGoogle), ic.Google)
		goth.UseProviders(google.New(id, secret, rt.MakeOAuthCallbackURL(baseURL, ProviderNameGoogle)))
	}
}
