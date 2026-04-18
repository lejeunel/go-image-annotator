package image

import (
	"fmt"
	"github.com/google/uuid"
	uuidw "github.com/lejeunel/go-image-annotator-v2/shared/uuid"
)

type ImageId struct {
	uuidw.UUIDWrapper[ImageId]
}

func NewImageId() ImageId {
	return ImageId{uuidw.UUIDWrapper[ImageId]{UUID: uuid.New()}}
}

func NewImageIdFromString(s string) (ImageId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ImageId{}, fmt.Errorf("invalid ImageId: %w", err)
	}

	return ImageId{
		UUIDWrapper: uuidw.FromUUID[ImageId](id),
	}, nil
}
