package list

import (
	m "github.com/lejeunel/go-image-annotator/entities/meta"
)

type OutputPort interface {
	Error(error)
	SuccessListMetadata([]m.MetaData)
}
