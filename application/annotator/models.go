package annotator

import (
	"io"

	an "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
)

type ImageInfo struct {
	Id         string
	Collection string
}

type Image struct {
	Reader      io.Reader
	Id          string
	Collection  string
	MIMEType    string
	Annotations Annotations
}

type Annotations struct {
	BoundingBoxes []an.BoundingBox
}

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

type BoundingBox struct {
	Id     string
	Label  string
	Color  string
	Xc     float32
	Yc     float32
	Width  float32
	Height float32
}

func NewImageInfo(imageId im.ImageId, collection string) ImageInfo {
	return ImageInfo{Id: imageId.String(), Collection: collection}
}

func NewImage(id im.ImageId, reader io.Reader, collection string, mimetype string,
) Image {
	return Image{Id: id.String(), Collection: collection, Reader: reader, MIMEType: mimetype}
}
