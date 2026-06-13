package server

import (
	"github.com/lejeunel/go-image-annotator/app"
	"log/slog"
)

type Server struct {
	*app.Interactors
	slog.Logger
}

func NewServer(interactors *app.Interactors, logger slog.Logger) *Server {
	return &Server{interactors, logger}
}
