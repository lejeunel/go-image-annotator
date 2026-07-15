package sqlite

import (
	a "github.com/lejeunel/go-image-annotator/app"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	ing "github.com/lejeunel/go-image-annotator/modules/ingester"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
	"github.com/lejeunel/go-image-annotator/modules/token"
)

func NewSQLiteInteractors(i SQLiteRepos, pageSize int, allowedImageFormats []string,
	APItokenGenerator token.Interface,
	passwordGenerator token.Interface,
	forgotPasswordTokenGen token.Interface,
	forgotPasswordTokenExpirationMinutes int,
	passwordValidator pv.MyPasswordValidator,
	passwordHasher token.TokenHasher,
	ingester ing.Interface,
	auth auth.Authorizer) a.Interactors {

	return a.Interactors{
		Label:      NewSQLiteLabelInteractors(i.Label, pageSize, auth),
		Collection: NewSQLiteCollectionInteractors(i.Collection, i.Group, pageSize, auth),
		Image:      NewSQLiteImageInteractors(i, ingester, pageSize, auth),
		User: NewSQLiteUserInteractors(i,
			APItokenGenerator,
			forgotPasswordTokenGen,
			passwordValidator,
			passwordHasher,
			forgotPasswordTokenExpirationMinutes,
			passwordGenerator, auth),
		Annotation: NewSQLiteAnnotationInteractors(i, auth),
		Group:      NewSQLiteGroupInteractors(i.Group),
	}
}
