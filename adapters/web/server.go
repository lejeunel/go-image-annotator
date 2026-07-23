package web

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	app "github.com/lejeunel/go-image-annotator/app"
	s "github.com/lejeunel/go-image-annotator/shared/session"
)

type Server struct {
	*app.Interactors
	b.PageBuilder
	s.SessionManager
	DefaultPageSize int
}

func NewServer(
	interactors *app.Interactors,
	pageBuilder b.PageBuilder,
	sessionManager s.SessionManager,
	pageSize int) *Server {
	return &Server{
		Interactors:     interactors,
		SessionManager:  sessionManager,
		PageBuilder:     pageBuilder,
		DefaultPageSize: pageSize,
	}
}
