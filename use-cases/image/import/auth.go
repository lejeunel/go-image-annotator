package import_image

import (
	"context"
)

type Auth interface {
	ImportImage(ctx context.Context, destinationGroup string) error
}
