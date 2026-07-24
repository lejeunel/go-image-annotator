package find

import (
	r "github.com/lejeunel/go-image-annotator/entities/role"
)

type OutputPort interface {
	Error(error)
	SuccessFindRole(r.Role)
}
