package server

import (
	"net/http"
)

func RegisterAPIEndpoints(mux *http.ServeMux, server Server, path string) {
	HandlerFromMuxWithBaseURL(&server, mux, path)
}
