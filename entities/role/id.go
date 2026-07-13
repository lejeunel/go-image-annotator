package role

import (
	"github.com/google/uuid"
	uuidw "github.com/lejeunel/go-image-annotator/shared/uuid"
)

type RoleId struct {
	uuidw.UUIDWrapper[RoleId]
}

func NewRoleId() RoleId {
	return RoleId{uuidw.UUIDWrapper[RoleId]{UUID: uuid.New()}}
}
