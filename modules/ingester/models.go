package ingester

import (
	"io"

	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Request struct {
	UserId        u.UserId
	Collection    string
	Labels        []string
	BoundingBoxes []an.BoundingBoxRequest
	Polygons      []an.PolygonRequest
	Reader        io.Reader
}

type Response struct {
	ImageId    im.ImageId
	Collection string
}

type BatchResponse struct {
	NumIngestedImages int64
}

type BatchRequest struct {
	UserId     u.UserId
	Collection string
	ReaderAt   io.ReaderAt
	Size       int64
}
