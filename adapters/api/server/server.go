package server

import (
	itrs "github.com/lejeunel/go-image-annotator/app/interactors"
	"log/slog"
)

type Server struct {
	*itrs.Interactors
	slog.Logger
}

func NewServer(interactors *itrs.Interactors, logger slog.Logger) *Server {
	return &Server{interactors, logger}
}
