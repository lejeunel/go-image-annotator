package role

import (
	"github.com/google/uuid"
	uuidw "github.com/lejeunel/go-image-annotator/shared/uuid"
)

type ProfileId struct {
	uuidw.UUIDWrapper[ProfileId]
}

func NewProfileId() ProfileId {
	return ProfileId{uuidw.UUIDWrapper[ProfileId]{UUID: uuid.New()}}
}
