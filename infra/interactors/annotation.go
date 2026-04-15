package interactors

import (
	infra "github.com/lejeunel/go-image-annotator-v2/infra"
	an "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate"
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
	remano "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
)

func NewSQLiteAnnotationInteractors(repos *infra.SQLiteInfra) *an.Interactors {
	return &an.Interactors{
		AddBox: *addbox.NewInteractor(repos.ImageStore, repos.AnnotationRepo),
		Delete: *remano.NewInteractor(repos.AnnotationRepo),
	}
}
