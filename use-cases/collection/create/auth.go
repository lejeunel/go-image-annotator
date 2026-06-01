package create

import (
	"github.com/lejeunel/go-image-annotator/shared/auth"
)

type Auth interface {
	CreateCollection(p auth.PrincipalProvider, group string) error
}
