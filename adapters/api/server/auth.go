package server

import (
	p "github.com/lejeunel/go-image-annotator/entities/principal"
	"net/http"
)

type HTTPPrincipalProvider interface {
	Provide(*http.Request) (*p.Principal, error)
}
