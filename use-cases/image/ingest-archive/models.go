package ingest

import (
	"io"

	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Request struct {
	Collection string
	Reader     io.Reader
}

type Response struct {
	Id     t.TaskId
	Issuer u.UserId
	Type   t.TaskType
}
