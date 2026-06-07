package web

import (
	a "github.com/lejeunel/go-image-annotator/app/annotator"
	"github.com/lejeunel/go-image-annotator/shared/html"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	s "github.com/lejeunel/go-image-annotator/shared/session"
	u "github.com/lejeunel/go-image-annotator/use-cases"
)

type Server struct {
	*u.Interactors
	html.PageBuilder
	annotator *a.Annotator
	s.SessionManager
	ip.OAuthHandler
}

func NewServer(interactors *u.Interactors, annotator *a.Annotator,
	pageBuilder html.PageBuilder, sessionManager s.SessionManager,
	identityHandler ip.OAuthHandler) *Server {
	return &Server{Interactors: interactors, annotator: annotator,
		SessionManager: sessionManager, PageBuilder: pageBuilder,
		OAuthHandler: identityHandler}
}
