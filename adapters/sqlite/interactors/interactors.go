package interactors

import (
	"github.com/lejeunel/go-image-annotator/adapters/sqlite"
	tok "github.com/lejeunel/go-image-annotator/app/token"
	u "github.com/lejeunel/go-image-annotator/use-cases"
)

func NewSQLiteInteractors(repos *infra.SQLiteInfra, pageSize int, allowedImageFormats []string,
	tokenGenerator tok.TokenGenerator) *u.Interactors {

	return &u.Interactors{
		Label:      NewSQLiteLabelInteractors(repos.Label, pageSize),
		Collection: NewSQLiteCollectionInteractors(repos.Collection, repos.Group, pageSize),
		Image:      NewSQLiteImageInteractors(repos, allowedImageFormats),
		User:       NewSQLiteUserInteractors(repos, tokenGenerator),
		Annotation: NewSQLiteAnnotationInteractors(repos),
	}
}
