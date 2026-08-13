package sqlite

import (
	anr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	imr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	lbr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	ims "github.com/lejeunel/go-image-annotator/modules/image-store"
	an "github.com/lejeunel/go-image-annotator/use-cases/annotate"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	updpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-polygon"
	remano "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

func NewAnnotationInteractors(ims ims.ImageStore,
	imr imr.ImageRepo,
	lbr lbr.LabelRepo,
	anr anr.AnnotationRepo,
	auth auth.Interface,
) an.Interactors {
	return an.Interactors{
		AddPolygon:    addpoly.New(ims, anr, lbr, addpoly.WithAuth(auth)),
		UpdatePolygon: updpoly.New(anr, lbr, updpoly.WithAuth(auth)),
		AddBox:        addbox.New(ims, anr, lbr, addbox.WithAuth(auth)),
		UpdateBox:     updbox.New(anr, lbr, updbox.WithAuth(auth)),
		Delete:        remano.New(anr, remano.WithAuth(auth)),
		UpdateLabel:   updlbl.New(anr, lbr, updlbl.WithAuth(auth)),
		AddImageLabel: addlbl.New(anr, lbr, ims, addlbl.WithAuth(auth)),
	}
}
