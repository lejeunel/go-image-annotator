package sqlite

import (
	a "github.com/lejeunel/go-image-annotator/app"
	tok "github.com/lejeunel/go-image-annotator/modules/token"
)

func NewSQLiteInteractors(i SQLiteRepos, pageSize int, allowedImageFormats []string,
	tokenGenerator tok.TokenGenerator) a.Interactors {

	return a.Interactors{
		Label:      NewSQLiteLabelInteractors(i.Label, pageSize),
		Collection: NewSQLiteCollectionInteractors(i.Collection, i.Group, pageSize),
		Image:      NewSQLiteImageInteractors(i, allowedImageFormats, pageSize),
		User:       NewSQLiteUserInteractors(i, tokenGenerator),
		Annotation: NewSQLiteAnnotationInteractors(i),
	}
}
