package web

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	app "github.com/lejeunel/go-image-annotator/app"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	s "github.com/lejeunel/go-image-annotator/shared/session"
)

type Server struct {
	*app.Interactors
	b.PageBuilder
	b.UserDashboardBuilder
	a.Annotator
	s.SessionManager
	ip.OAuthHandler
}

func NewServer(interactors *app.Interactors, annotator a.Annotator,
	pageBuilder b.PageBuilder, sessionManager s.SessionManager,
	identityHandler ip.OAuthHandler) *Server {
	return &Server{Interactors: interactors, Annotator: annotator,
		SessionManager: sessionManager, PageBuilder: pageBuilder,
		UserDashboardBuilder: b.NewUserDashboardBuilder(),
		OAuthHandler:         identityHandler}
}
