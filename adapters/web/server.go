package web

import (
	ap "github.com/lejeunel/go-image-annotator/adapters/web/annotator/presenters"
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
	ip.AuthHandler
	ap.AnnotationPagePresenter
	ap.AnnotoriousPresenter
	DefaultPageSize int
}

func NewServer(
	interactors *app.Interactors,
	annotator a.Annotator,
	pageBuilder b.PageBuilder,
	annotationPagePresenter ap.AnnotationPagePresenter,
	annotoriousPresenter ap.AnnotoriousPresenter,
	sessionManager s.SessionManager,
	identityHandler ip.AuthHandler,
	pageSize int) *Server {
	return &Server{
		Interactors:             interactors,
		Annotator:               annotator,
		SessionManager:          sessionManager,
		PageBuilder:             pageBuilder,
		UserDashboardBuilder:    b.NewUserDashboardBuilder(),
		AnnotationPagePresenter: annotationPagePresenter,
		AnnotoriousPresenter:    annotoriousPresenter,
		AuthHandler:             identityHandler,
		DefaultPageSize:         pageSize,
	}
}
