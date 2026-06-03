package import_image

import (
	"context"
)

type Auth interface {
	ImportImage(ctx context.Context, sourceGroup string, destinationGroup string) error
}
