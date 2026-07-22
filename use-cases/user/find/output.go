package find

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type OutputPort interface {
	SuccessFindUser(u.User)
	Error(error)
}
