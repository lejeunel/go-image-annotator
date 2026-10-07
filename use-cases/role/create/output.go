package create

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
)

type OutputPort interface {
	SuccessCreateRole(rl.Role)
	Error(error)
}
