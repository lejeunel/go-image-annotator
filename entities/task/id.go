package task

import (
	"fmt"
	"github.com/google/uuid"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	uuidw "github.com/lejeunel/go-image-annotator/shared/uuid"
)

type TaskId struct{ uuidw.UUIDWrapper[TaskId] }

func NewTaskId() TaskId {
	return TaskId{uuidw.UUIDWrapper[TaskId]{UUID: uuid.New()}}
}

func NewTaskIdFromString(s string) (*TaskId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("invalid TaskId: %w: %w", err, e.ErrValidation)
	}

	return &TaskId{
		UUIDWrapper: uuidw.FromUUID[TaskId](id),
	}, nil
}
