package auth

import (
	"log/slog"
	"os"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	rt "github.com/lejeunel/go-image-annotator/routes"
	sm "github.com/lejeunel/go-image-annotator/shared/session"
	reqpw "github.com/lejeunel/go-image-annotator/use-cases/user/forgot-password"
	respw "github.com/lejeunel/go-image-annotator/use-cases/user/reset-forgotten-password"
)

type Server struct {
	b.LoginPageBuilder
	b.ForgotPasswordBuilder
	b.ResetPasswordBuilder
	slog.Logger
	sm.SessionManager
	requestTokenItr  reqpw.Interactor
	ResetPasswordItr respw.Interactor
	baseURL          string
}

func NewAuthWebServer(
	baseURL string,
	basePageBuilder b.BasePageBuilder,
	logger slog.Logger,
	sm sm.SessionManager,
	reqForgottenPw reqpw.Interactor,
	resForgottenPw respw.Interactor,
) Server {

	loginPageBuilder := b.NewLoginPageBuilder(basePageBuilder)
	loginPageBuilder.AddOAuthProvider("google", rt.MakeOAuthLoginURL("google"))
	SetupForGoogle(OAuthProviderConfig{Key: os.Getenv("GOIA_GOOGLE_CLIENT_ID"),
		Secret:      os.Getenv("GOIA_GOOGLE_CLIENT_SECRET"),
		CallbackURL: rt.MakeOAuthCallbackURL(baseURL, "google")})
	return Server{
		loginPageBuilder,
		b.NewForgotPasswordBuilder(basePageBuilder),
		b.NewResetPasswordBuilder(basePageBuilder),
		logger,
		sm,
		reqForgottenPw,
		resForgottenPw,
		baseURL,
	}
}
