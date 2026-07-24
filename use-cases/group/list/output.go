package list

import (
	g "github.com/lejeunel/go-image-annotator/entities/group"
)

type OutputPort interface {
	SuccessListGroups([]g.Group)
	Error(error)
}
