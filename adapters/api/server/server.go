package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	itrs "github.com/lejeunel/go-image-annotator/app/interactors"
)

type Server struct {
	*itrs.Interactors
	slog.Logger
}

func NewServer(interactors *itrs.Interactors, logger slog.Logger) *Server {
	return &Server{interactors, logger}
}

// NewParamErrorHandler builds a handler that catches malformed requests.
func NewParamErrorHandler(l slog.Logger) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		l.Warn("rejecting malformed request", "path", r.URL.Path, "error", err)
		json.WriteError(w, http.StatusBadRequest, makeParamErrorMessage(err))
	}
}

func makeParamErrorMessage(err error) string {
	var invalidFormat *InvalidParamFormatError
	if errors.As(err, &invalidFormat) {
		return fmt.Sprintf("invalid value for parameter %q", invalidFormat.ParamName)
	}

	var required *RequiredParamError
	if errors.As(err, &required) {
		return fmt.Sprintf("missing required parameter %q", required.ParamName)
	}

	return "invalid request parameters"
}
