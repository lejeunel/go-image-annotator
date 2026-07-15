package clone

import (
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Request struct {
	Source           string
	Destination      string
	DestinationGroup *string
	Deep             bool
}

type Response struct {
	Id     t.TaskId
	Issuer u.UserId
	Type   t.TaskType
	State  t.TaskState
}
