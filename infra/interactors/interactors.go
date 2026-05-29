package interactors

import (
	"github.com/lejeunel/go-image-annotator/infra"
	u "github.com/lejeunel/go-image-annotator/use-cases"
)

func NewSQLiteInteractors(repos *infra.SQLiteInfra, pageSize int, allowedImageFormats []string) *u.Interactors {

	return &u.Interactors{
		Label:      NewSQLiteLabelInteractors(repos.LabelRepo, pageSize),
		Collection: NewSQLiteCollectionInteractors(repos.CollectionRepo, pageSize),
		Image:      NewSQLiteImageInteractors(repos, allowedImageFormats),
		Annotation: NewSQLiteAnnotationInteractors(repos),
	}
}
