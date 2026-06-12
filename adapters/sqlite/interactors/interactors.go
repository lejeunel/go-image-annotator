package interactors

import (
	i "github.com/lejeunel/go-image-annotator/adapters/sqlite/infra"
	tok "github.com/lejeunel/go-image-annotator/app/token"
	u "github.com/lejeunel/go-image-annotator/use-cases"
)

func NewSQLiteInteractors(i i.SQLiteInfra, pageSize int, allowedImageFormats []string,
	tokenGenerator tok.TokenGenerator) *u.Interactors {

	return &u.Interactors{
		Label:      NewSQLiteLabelInteractors(i.Label, pageSize),
		Collection: NewSQLiteCollectionInteractors(i.Collection, i.Group, pageSize),
		Image:      NewSQLiteImageInteractors(i, allowedImageFormats),
		User:       NewSQLiteUserInteractors(i, tokenGenerator),
		Annotation: NewSQLiteAnnotationInteractors(i),
	}
}
