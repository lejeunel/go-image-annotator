package delete

import (
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Response struct {
	Id     t.TaskId
	Issuer u.UserId
	Type   t.TaskType
}
