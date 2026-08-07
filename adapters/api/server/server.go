package server

import (
	"log/slog"

	itrs "github.com/lejeunel/go-image-annotator/app/interactors"
)

type Server struct {
	*itrs.Interactors
	slog.Logger
}

func NewServer(interactors *itrs.Interactors, logger slog.Logger) *Server {
	return &Server{interactors, logger}
}
