package sqlite

import (
	an "github.com/lejeunel/go-image-annotator/use-cases/annotate"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	remano "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

func NewSQLiteAnnotationInteractors(repos SQLiteRepos) an.Interactors {
	return an.Interactors{
		AddPolygon:    addpoly.New(repos.ImageStore, repos.Annotation, repos.Label),
		AddBox:        addbox.New(repos.ImageStore, repos.Annotation, repos.Label),
		UpdateBox:     updbox.New(repos.Annotation, repos.Label),
		Delete:        remano.New(repos.Annotation),
		UpdateLabel:   updlbl.New(repos.Annotation, repos.Label),
		AddImageLabel: addlbl.New(repos.Annotation, repos.Label, repos.ImageStore),
	}
}
