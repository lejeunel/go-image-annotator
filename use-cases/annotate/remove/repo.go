package remove

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
)

type Repo interface {
	RemoveAnnotation(a.AnnotationId) error
}
