package raw

import (
	"io"
)

type OutputPort interface {
	SuccessReadRawImage(io.Reader)
	Error(error)
}
