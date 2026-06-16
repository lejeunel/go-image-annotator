package view

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
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

type Annotation struct {
	Id    string
	Label string
}

type BoundingBox struct {
	Id     string
	Label  string
	Color  string
	Xc     float32
	Yc     float32
	Width  float32
	Height float32
	Angle  float32
	Author string
	Time   string
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
