package form

import (
	"io"
)

type Renderer interface {
	Render(io.Writer)
}
