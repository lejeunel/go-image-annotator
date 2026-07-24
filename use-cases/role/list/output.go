package list

import (
	r "github.com/lejeunel/go-image-annotator/entities/role"
)

type OutputPort interface {
	SuccessListRoles([]r.Role)
	Error(error)
}
