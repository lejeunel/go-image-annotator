package interactors

import (
	"github.com/lejeunel/go-image-annotator/adapters/sqlite"
	tg "github.com/lejeunel/go-image-annotator/app/token-generator"
	u "github.com/lejeunel/go-image-annotator/use-cases"
)

func NewSQLiteInteractors(repos *infra.SQLiteInfra, pageSize int, allowedImageFormats []string) *u.Interactors {

	return &u.Interactors{
		Label:      NewSQLiteLabelInteractors(repos.Label, pageSize),
		Collection: NewSQLiteCollectionInteractors(repos.Collection, repos.Group, pageSize),
		Image:      NewSQLiteImageInteractors(repos, allowedImageFormats),
		User:       NewSQLiteUserInteractors(repos, tg.NewTokenGenerator(32)),
		Annotation: NewSQLiteAnnotationInteractors(repos),
	}
}
