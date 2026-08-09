package ingest

import (
	ing "github.com/lejeunel/go-image-annotator/modules/image-ingester"
)

type OutputPort interface {
	Success(ing.Response)
	Error(error)
}
