package group

import (
	"github.com/google/uuid"
	uuidw "github.com/lejeunel/go-image-annotator/shared/uuid"
)

type GroupId struct {
	uuidw.UUIDWrapper[GroupId]
}

func NewGroupId() GroupId {
	return GroupId{uuidw.UUIDWrapper[GroupId]{UUID: uuid.New()}}
}
