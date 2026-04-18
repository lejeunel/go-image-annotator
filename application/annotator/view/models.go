package view

import (
	an "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	"io"
)

type ScrollerButton struct {
	IsActive   bool
	Text       string
	ImageId    string
	Collection string
}

type ScrollerButtons struct {
	Next ScrollerButton
	Prev ScrollerButton
}

type ImageLabel struct {
	Id    string
	Label string
}

type ImageInfo struct {
	Id         string
	Collection string
	Specs      im.ImageSpecs
}

type BoundingBox struct {
	Id     string
	Label  string
	Color  string
	Xc     float32
	Yc     float32
	Width  float32
	Height float32
}

type Annotations struct {
	BoundingBoxes []BoundingBox
}

type Image struct {
	Reader      io.Reader
	Id          string
	Collection  string
	MIMEType    string
	Annotations Annotations
}

func NewImageInfo(imageId im.ImageId, collection string, specs im.ImageSpecs) ImageInfo {
	return ImageInfo{Id: imageId.String(), Collection: collection, Specs: specs}
}

func NewImage(id im.ImageId, reader io.Reader, collection string, mimetype string,
) Image {
	return Image{Id: id.String(), Collection: collection, Reader: reader, MIMEType: mimetype}
}

func MakeBoundingBox(b *an.BoundingBox, colorIndex int) *BoundingBox {
	return &BoundingBox{
		Id:     b.Id.String(),
		Label:  b.Label.Name,
		Color:  Palette[colorIndex%(len(Palette)-1)],
		Xc:     b.Xc,
		Yc:     b.Yc,
		Width:  b.Width,
		Height: b.Height,
	}
}

func MakeImageLabels(labels []*an.ImageLabel) []*ImageLabel {
	result := []*ImageLabel{}
	for _, l := range labels {
		result = append(result, &ImageLabel{Id: l.Id.String(),
			Label: l.Label.Name})
	}
	return result
}

func MakeBoundingBoxes(boxes []*an.BoundingBox) []*BoundingBox {
	result := []*BoundingBox{}
	for i, b := range boxes {
		result = append(result, MakeBoundingBox(b, i))
	}
	return result
}
