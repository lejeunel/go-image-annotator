package auth

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func (s Server) OAuthLogin(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (s Server) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		s.Logger.Error("completing oauth login", "error", err)
		http.Error(w, "oauth login failed", http.StatusUnauthorized)
		return
	}

	if err := s.SessionManager.FinishOAuthLogin(r.Context(), user.Email); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// MinSessionSecretLength is the shortest secret we accept for signing the
// cookie.
const MinSessionSecretLength = 32

// SetupOAuthSessionStore configures an in-memory cookie store.
func SetupOAuthSessionStore(secret string, logger slog.Logger) {
	key := []byte(secret)

	if len(goth.GetProviders()) > 0 {
		if len(key) < MinSessionSecretLength {
			logger.Error(
				"an oauth provider is configured but GOIA_SESSION_SECRET is missing or too short",
				"required_bytes", MinSessionSecretLength, "got_bytes", len(key))
			os.Exit(1)
		}
	}

	store := sessions.NewCookieStore(key)
	store.Options.HttpOnly = true
	gothic.Store = store
}

func MaybeSetupGoogle(
	pb *b.LoginPageBuilder,
	baseURL string,
	logger slog.Logger,
	id, secret string,
) {
	if (id != "") && (secret != "") {
		logger.Info("setting up google auth")
		pb.AddOAuthProvider(ProviderNameGoogle, rt.MakeOAuthLoginURL(ProviderNameGoogle), ic.Google)
		goth.UseProviders(
			google.New(id, secret, rt.MakeOAuthCallbackURL(baseURL, ProviderNameGoogle)),
		)
	}
}

func MaybeSetupGithub(
	pb *b.LoginPageBuilder,
	baseURL string,
	logger slog.Logger,
	id, secret string,
) {
	if (id != "") && (secret != "") {
		logger.Info("setting up github auth")
		pb.AddOAuthProvider(ProviderNameGithub, rt.MakeOAuthLoginURL(ProviderNameGithub), ic.Github)
		goth.UseProviders(
			github.New(
				id,
				secret,
				rt.MakeOAuthCallbackURL(baseURL, ProviderNameGithub),
				"user:email",
			),
		)
	}
}
