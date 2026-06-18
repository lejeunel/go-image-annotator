package sqlite

import (
	a "github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/modules/auth"
	tok "github.com/lejeunel/go-image-annotator/modules/token"
)

func NewSQLiteInteractors(i SQLiteRepos, pageSize int, allowedImageFormats []string,
	tokenGenerator tok.TokenGenerator, auth auth.Auth) a.Interactors {

	return a.Interactors{
		Label:      NewSQLiteLabelInteractors(i.Label, pageSize, auth),
		Collection: NewSQLiteCollectionInteractors(i.Collection, i.Group, pageSize, auth),
		Image:      NewSQLiteImageInteractors(i, allowedImageFormats, pageSize, auth),
		User:       NewSQLiteUserInteractors(i, tokenGenerator, auth),
		Annotation: NewSQLiteAnnotationInteractors(i, auth),
	}
}
