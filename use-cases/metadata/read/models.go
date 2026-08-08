package read

import (
	m "github.com/lejeunel/go-image-annotator/entities/meta"
)

type Request struct {
	ImageId    string
	Collection string
	Key        string
}

type Response struct {
	ImageId    string
	Collection string
	Data       m.MetaData
}
