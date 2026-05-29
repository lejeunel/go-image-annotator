package label

import (
	"github.com/google/uuid"
	uuidw "github.com/lejeunel/go-image-annotator/shared/uuid"
)

type LabelId struct{ uuidw.UUIDWrapper[LabelId] }

func NewLabelId() LabelId {
	return LabelId{uuidw.UUIDWrapper[LabelId]{UUID: uuid.New()}}
}
