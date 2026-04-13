package annotator

import (
	"io"

	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
)

type ImageInfo struct {
	Id         string
	Collection string
}

type Image struct {
	Reader     io.Reader
	Id         string
	Collection string
	MIMEType   string
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

func NewImageInfo(imageId im.ImageId, collection string) ImageInfo {
	return ImageInfo{Id: imageId.String(), Collection: collection}
}

func NewImage(id im.ImageId, reader io.Reader, collection string, mimetype string) Image {
	return Image{Id: id.String(), Collection: collection, Reader: reader, MIMEType: mimetype}
}
