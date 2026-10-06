package auth

import (
	"log/slog"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cfg "github.com/lejeunel/go-image-annotator/config"
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
	Notifier
}

func New(
	baseURL string,
	cfg cfg.Config,
	basePageBuilder b.BasePageBuilder,
	logger slog.Logger,
	sm sm.SessionManager,
	nf Notifier,
	reqForgottenPw reqpw.Interactor,
	resForgottenPw respw.Interactor,
) Server {
	loginPageBuilder := b.NewLoginPageBuilder(basePageBuilder)
	MaybeSetupGoogle(&loginPageBuilder, baseURL, logger, cfg.GoogleClientId, cfg.GoogleClientSecret)
	MaybeSetupGithub(&loginPageBuilder, baseURL, logger, cfg.GithubClientId, cfg.GithubClientSecret)
	SetupOAuthSessionStore(cfg.SessionSecret, logger)
	return Server{
		loginPageBuilder,
		b.NewForgotPasswordBuilder(basePageBuilder),
		b.NewResetPasswordBuilder(basePageBuilder),
		logger,
		sm,
		reqForgottenPw,
		resForgottenPw,
		baseURL,
		nf,
	}
}
