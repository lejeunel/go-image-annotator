package image

import (
	"fmt"

	"github.com/google/uuid"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	uuidw "github.com/lejeunel/go-image-annotator/shared/uuid"
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
		return ImageId{}, fmt.Errorf("invalid ImageId: %w: %w", err, e.ErrValidation)
	}

	return ImageId{
		UUIDWrapper: uuidw.FromUUID[ImageId](id),
	}, nil
}
