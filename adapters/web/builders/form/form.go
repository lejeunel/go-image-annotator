package form

import (
	. "maragu.dev/gomponents"
)

type FormField interface {
	Build() Node
}
