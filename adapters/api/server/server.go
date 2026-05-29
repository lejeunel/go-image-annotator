package server

import (
	u "github.com/lejeunel/go-image-annotator/use-cases"
)

type Server struct {
	*u.Interactors
}

func NewServer(interactors *u.Interactors) *Server {
	return &Server{interactors}
}
