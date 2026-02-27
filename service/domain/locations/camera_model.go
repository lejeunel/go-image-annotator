package locations

import (
	e "datahub/errors"
	g "datahub/generic"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"time"
)

type CameraId struct{ g.UUIDWrapper[CameraId] }

func NewCameraId() *CameraId {
	id := uuid.New()
	return &CameraId{g.UUIDWrapper[CameraId]{UUID: id}}

}

func (id *CameraId) Equal(to *CameraId) bool {
	if (to == nil) && (id == nil) {
		return true
	}
	if id.String() == to.String() {
		return true
	}
	return false
}

func NewCameraIdFromUUID(id uuid.UUID) *CameraId {
	return &CameraId{g.UUIDWrapper[CameraId]{UUID: id}}
}

func NewCameraIdFromString(s string) (*CameraId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parsing camera id: %w", e.ErrValidation)
	}
	return NewCameraIdFromUUID(id), nil

}

type CameraOption func(*Camera)

func WithTransmitter(transmitter string) CameraOption {
	return func(c *Camera) {
		c.Transmitter = transmitter
	}
}
func NewCamera(name string, site *Site, opts ...CameraOption) (*Camera, error) {
	cam := &Camera{Id: *NewCameraId(),
		Name:  name,
		Site:  site,
		Group: site.Group}

	for _, opt := range opts {
		opt(cam)
	}
	if err := cam.Validate(); err != nil {
		return nil, err
	}
	return cam, nil
}

type Camera struct {
	Id          CameraId
	Name        string
	Site        *Site
	Transmitter string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Group       string
}

func (c Camera) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Name, validation.Match(g.ResourceNameRegExp)),
	); err != nil {
		return fmt.Errorf("validating camera name (%v): %w", c.Name, e.ErrResourceName)
	}
	return nil
}
