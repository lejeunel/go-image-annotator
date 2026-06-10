package server

import (
	u "github.com/lejeunel/go-image-annotator/use-cases"
	"log/slog"
)

type Server struct {
	*u.Interactors
	slog.Logger
}

func NewServer(interactors *u.Interactors, logger slog.Logger) *Server {
	return &Server{interactors, logger}
}
