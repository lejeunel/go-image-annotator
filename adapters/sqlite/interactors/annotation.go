package interactors

import (
	infra "github.com/lejeunel/go-image-annotator/adapters/infra"
	an "github.com/lejeunel/go-image-annotator/use-cases/annotate"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	remano "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

func NewSQLiteAnnotationInteractors(repos *infra.SQLiteInfra) *an.Interactors {
	return &an.Interactors{
		AddBox:        *addbox.NewInteractor(repos.ImageStore, repos.AnnotationRepo),
		UpdateBox:     *updbox.NewInteractor(repos.AnnotationRepo),
		Delete:        *remano.NewInteractor(repos.AnnotationRepo),
		UpdateLabel:   *updlbl.NewInteractor(repos.AnnotationRepo),
		AddImageLabel: *addlbl.NewInteractor(repos.AnnotationRepo, repos.ImageStore),
	}
}
