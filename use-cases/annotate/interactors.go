package collection

import (
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/update-label"
)

type Interactors struct {
	AddBox      addbox.Interactor
	UpdateBox   updbox.Interactor
	Delete      remove.Interactor
	UpdateLabel updlbl.Interactor
}
