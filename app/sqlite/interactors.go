package sqlite

import (
	a "github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/modules/auth"
	au "github.com/lejeunel/go-image-annotator/modules/authentifier"
)

func NewSQLiteInteractors(i SQLiteRepos, pageSize int, allowedImageFormats []string,
	tokenGenerator au.AuthGenerator, passwordGenerator au.AuthGenerator, auth auth.Auth) a.Interactors {

	return a.Interactors{
		Label:      NewSQLiteLabelInteractors(i.Label, pageSize, auth),
		Collection: NewSQLiteCollectionInteractors(i.Collection, i.Group, pageSize, auth),
		Image:      NewSQLiteImageInteractors(i, allowedImageFormats, pageSize, auth),
		User:       NewSQLiteUserInteractors(i, tokenGenerator, passwordGenerator, auth),
		Annotation: NewSQLiteAnnotationInteractors(i, auth),
	}
}
