package delete

import (
	"github.com/lejeunel/go-image-annotator/shared/auth"
)

type Auth interface {
	DeleteCollection(p auth.PrincipalProvider, group string) error
}
