package scroll

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Request struct {
	CurrentImageId    string
	CurrentCollection clc.CollectionName
	im.FilterStr
	im.OrderStr
}

type Response struct {
	Next *im.BaseImage
	Prev *im.BaseImage
	im.FilterStr
	im.OrderStr
}
