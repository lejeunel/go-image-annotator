package ingester

import (
	"io"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Response struct {
	ImageIds   []im.ImageId
	Collection string
}

type Request struct {
	UserId     u.UserId
	Collection string
	ReaderAt   io.ReaderAt
	Size       int64
}
