package images

import (
	"context"
	"io"
)

type ImageSetDescriber interface {
	Describe(ctx context.Context, w io.Writer) error
}
