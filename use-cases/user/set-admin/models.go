package set_admin

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Response struct {
	Id      u.UserId
	IsAdmin bool
}
