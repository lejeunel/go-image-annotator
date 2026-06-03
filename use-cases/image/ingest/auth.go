package ingest

import (
	"context"
)

type Auth interface {
	IngestImage(ctx context.Context, group string) error
}
