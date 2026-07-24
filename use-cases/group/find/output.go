package find

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
)

type OutputPort interface {
	Error(error)
	SuccessFindGroup(grp.Group)
}
