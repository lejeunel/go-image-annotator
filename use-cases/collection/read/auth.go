package read

import (
	"context"
)

type Auth interface {
	ReadCollection(context.Context) error
}
