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
	Reader        io.Reader
}

type Response struct {
	ImageId    im.ImageId
	Collection string
}
