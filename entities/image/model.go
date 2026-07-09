package image

import (
	"fmt"
	"io"
	"time"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type RawImage struct {
	Data []byte
	Hash string
}

type BaseImage struct {
	ImageId    ImageId
	Collection string
}

type ImageSpecs struct {
	MIMEType   string
	Width      int
	Height     int
	IngestedAt time.Time
}

type Image struct {
	Id            ImageId
	Collection    clc.Collection
	Specs         ImageSpecs
	Labels        []an.ImageLabel
	BoundingBoxes []an.BoundingBox
	Polygons      []an.Polygon
	Reader        io.Reader
	Hash          string
	MIMEType      string
}

func (i *Image) AddLabel(l lbl.Label) error {
	for _, label := range i.Labels {
		if l.Name == label.Label.Name {
			return fmt.Errorf("adding label %v to image: found duplicate: %w", label.Label.Name, e.ErrValidation)
		}
	}
	i.Labels = append(i.Labels, an.NewImageLabel(l))
	return nil
}
func (i *Image) AddBoundingBox(box a.BoundingBox) error {
	if err := a.ValidateBoundingBox(box.Xc, box.Yc, box.Width, box.Height, box.Angle); err != nil {
		return fmt.Errorf("adding bounding box to image: %w", err)
	}
	i.BoundingBoxes = append(i.BoundingBoxes, box)
	return nil
}
func (i *Image) NumAnnotations() int {
	return len(i.Labels) + len(i.BoundingBoxes) + len(i.Polygons)
}
func (i *Image) LabelNames() []string {
	labelNames := []string{}
	for _, label := range i.Labels {
		labelNames = append(labelNames, label.Label.Name)
	}
	return labelNames
}
func (i *Image) BoundingBoxSummary() []an.BoundingBoxResponse {
	summary := []an.BoundingBoxResponse{}
	for _, bbox := range i.BoundingBoxes {
		summary = append(summary,
			an.BoundingBoxResponse{Label: bbox.Label.Name, Xc: bbox.Xc, Yc: bbox.Yc, Width: bbox.Width, Height: bbox.Height})
	}
	return summary
}

func NewImage(id ImageId, collection clc.Collection) Image {
	return Image{Id: id, Collection: collection}
}
