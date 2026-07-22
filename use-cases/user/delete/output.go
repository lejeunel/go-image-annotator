package delete

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type OutputPort interface {
	Error(error)
	SuccessDeleteUser(u.UserId)
}
