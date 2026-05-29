package annotation

import (
	"fmt"
	"github.com/google/uuid"
	uuidw "github.com/lejeunel/go-image-annotator/shared/uuid"
)

type AnnotationId struct {
	uuidw.UUIDWrapper[AnnotationId]
}

func NewAnnotationId() AnnotationId {
	return AnnotationId{uuidw.UUIDWrapper[AnnotationId]{UUID: uuid.New()}}
}

func NewAnnotationIdFromString(s string) (*AnnotationId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("invalid ImageId: %w", err)
	}

	return &AnnotationId{
		UUIDWrapper: uuidw.FromUUID[AnnotationId](id),
	}, nil
}
