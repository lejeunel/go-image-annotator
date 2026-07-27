package find

import (
	"io"
)

type OutputPort interface {
	Error(error)
	SuccessReadPolicy(io.Reader)
}
