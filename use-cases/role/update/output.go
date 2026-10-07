package update

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
)

type OutputPort interface {
	SuccessUpdateRole(rl.Role)
	Error(error)
}
