package web

import (
	a "github.com/lejeunel/go-image-annotator/app/annotator"
	u "github.com/lejeunel/go-image-annotator/use-cases"
)

type Server struct {
	*u.Interactors
	annotator *a.Annotator
}

func NewServer(interactors *u.Interactors, annotator *a.Annotator) *Server {
	return &Server{Interactors: interactors, annotator: annotator}
}
