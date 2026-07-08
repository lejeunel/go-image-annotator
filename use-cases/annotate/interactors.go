package annotate

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	updpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-polygon"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

type Interactors struct {
	AddBox        addbox.Interactor
	UpdateBox     updbox.Interactor
	AddPolygon    addpoly.Interactor
	UpdatePolygon updpoly.Interactor
	Delete        remove.Interactor
	UpdateLabel   updlbl.Interactor
	AddImageLabel addlbl.Interactor
	Authorizer    auth.Authorizer
}
