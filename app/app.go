package app

import (
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	i "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	s "github.com/lejeunel/go-image-annotator/shared/session"
)

type App struct {
	Itrs           Interactors
	SessionManager s.MySessionManager
	i.OAuthHandler
	a.Annotator
}

func NewApp(itrs Interactors, sm s.MySessionManager, ip i.OAuthHandler, an a.Annotator) App {
	return App{itrs, sm, ip, an}
}
