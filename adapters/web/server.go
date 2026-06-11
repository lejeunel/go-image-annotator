package web

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	a "github.com/lejeunel/go-image-annotator/app/annotator"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	s "github.com/lejeunel/go-image-annotator/shared/session"
	u "github.com/lejeunel/go-image-annotator/use-cases"
)

type Server struct {
	*u.Interactors
	b.PageBuilder
	b.UserDashboardBuilder
	annotator *a.Annotator
	s.SessionManager
	ip.OAuthHandler
}

func NewServer(interactors *u.Interactors, annotator *a.Annotator,
	pageBuilder b.PageBuilder, sessionManager s.SessionManager,
	identityHandler ip.OAuthHandler) *Server {
	return &Server{Interactors: interactors, annotator: annotator,
		SessionManager: sessionManager, PageBuilder: pageBuilder,
		UserDashboardBuilder: b.NewUserDashboardBuilder(),
		OAuthHandler:         identityHandler}
}
