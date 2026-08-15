package scroll

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Request struct {
	CurrentImageId string
	im.FilterStr
	im.OrderStr
}

type Response struct {
	Next *im.BaseImage
	Prev *im.BaseImage
}
