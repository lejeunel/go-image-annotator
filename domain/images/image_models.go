package images

import (
	"bytes"
	"crypto/sha256"
	clc "datahub/domain/collections"
	loc "datahub/domain/locations"
	e "datahub/errors"
	g "datahub/generic"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	i "image"
	_ "image/jpeg"
	_ "image/png"
	"net/url"
	"time"
)

type ImageId struct{ g.UUIDWrapper[ImageId] }

func NewImageId() *ImageId {
	id := uuid.New()
	return &ImageId{g.UUIDWrapper[ImageId]{UUID: id}}

}

func (id *ImageId) Equal(to *ImageId) bool {
	if (to == nil) && (id == nil) {
		return true
	}
	if id.String() == to.String() {
		return true
	}
	return false
}

func NewImageIdFromUUID(id uuid.UUID) *ImageId {
	return &ImageId{g.UUIDWrapper[ImageId]{UUID: id}}
}

func NewImageIdFromString(s string) (*ImageId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parsing image id from string %v: %w", s, e.ErrValidation)
	}
	return NewImageIdFromUUID(id), nil

}

type ImageUpdatables struct {
	CapturedAt string `json:"captured_at"`
	Site       string `json:"site"`
	Camera     string `json:"camera"`
	Type_      string `json:"type"`
}

type BaseImage struct {
	Id           ImageId
	FileName     string
	CameraId     *loc.CameraId
	CapturedAt   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	SHA256       string
	MIMEType     string
	Width        int
	Height       int
	Type         string
	CollectionId clc.CollectionId
	Group        string
	Uri          url.URL
	Camera       *loc.Camera
}

func (im *BaseImage) GetSiteName() string {
	if im.Camera == nil {
		return ""
	}
	return im.Camera.Site.Name
}
func (im *BaseImage) GetCameraName() string {
	if im.Camera == nil {
		return ""
	}
	return im.Camera.Name
}

func (im *BaseImage) GetTransmitter() string {
	if im.Camera == nil {
		return ""
	}
	return im.Camera.Transmitter
}

type Image struct {
	Id            ImageId
	FileName      string
	CameraId      *loc.CameraId
	CapturedAt    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	SHA256        string
	MIMEType      string
	Width         int
	Height        int
	Type          string
	Group         string
	Data          []byte
	Annotations   []*Annotation
	BoundingBoxes []*BoundingBox
	Uri           url.URL
	Camera        *loc.Camera
	Collection    *clc.Collection
}

func New(data []byte) (*Image, error) {
	image := &Image{Data: data}
	h := sha256.New()
	h.Write(data)
	image.SHA256 = hex.EncodeToString(h.Sum(nil))
	im, format, err := i.Decode(bytes.NewBuffer(image.Data))
	if err != nil {
		return nil, err
	}

	image.Id = *NewImageId()
	image.Width = im.Bounds().Dx()
	image.Height = im.Bounds().Dy()
	image.MIMEType = "image/" + format

	return image, nil

}

func (im *Image) GetSiteName() string {
	if im.Camera == nil {
		return ""
	}
	return im.Camera.Site.Name
}
func (im *Image) GetCameraName() string {
	if im.Camera == nil {
		return ""
	}
	return im.Camera.Name
}

func (im *Image) GetTransmitter() string {
	if im.Camera == nil {
		return ""
	}
	return im.Camera.Transmitter
}

type ImageAnnotation struct {
	Id          uuid.UUID `db:"id"`
	LabelId     string    `db:"label_id"`
	ImageId     string    `db:"image_id"`
	AuthorEmail string    `db:"author_email"`
	CreatedAt   string    `db:"created_at"`
}

type LocationSelectorData struct {
	g.SelectorData
	ImageId      string
	CollectionId string
	Field        string
}
