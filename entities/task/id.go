package task

import (
	"github.com/google/uuid"
	uuidw "github.com/lejeunel/go-image-annotator/shared/uuid"
)

type TaskId struct{ uuidw.UUIDWrapper[TaskId] }

func NewTaskId() TaskId {
	return TaskId{uuidw.UUIDWrapper[TaskId]{UUID: uuid.New()}}
}
