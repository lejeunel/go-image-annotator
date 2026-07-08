package sqlite

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	an "github.com/lejeunel/go-image-annotator/use-cases/annotate"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	updpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-polygon"
	remano "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

func NewSQLiteAnnotationInteractors(repos SQLiteRepos, auth auth.Authorizer) an.Interactors {
	return an.Interactors{
		AddPolygon:    addpoly.New(repos.ImageStore, repos.Annotation, repos.Label, addpoly.WithAuth(auth)),
		UpdatePolygon: updpoly.New(repos.Annotation, repos.Label, updpoly.WithAuth(auth)),
		AddBox:        addbox.New(repos.ImageStore, repos.Annotation, repos.Label, addbox.WithAuth(auth)),
		UpdateBox:     updbox.New(repos.Annotation, repos.Label, updbox.WithAuth(auth)),
		Delete:        remano.New(repos.Annotation, remano.WithAuth(auth)),
		UpdateLabel:   updlbl.New(repos.Annotation, repos.Label, updlbl.WithAuth(auth)),
		AddImageLabel: addlbl.New(repos.Annotation, repos.Label, repos.ImageStore, addlbl.WithAuth(auth)),
	}
}
