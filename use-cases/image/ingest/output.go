package ingest

import (
	ing "github.com/lejeunel/go-image-annotator/modules/ingester"
)

type OutputPort interface {
	Success(ing.Response)
	Error(error)
}
