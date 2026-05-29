package interactors

import (
	an "github.com/lejeunel/go-image-annotator/use-cases/annotate"
	clc "github.com/lejeunel/go-image-annotator/use-cases/collection"
	im "github.com/lejeunel/go-image-annotator/use-cases/image"
	lbl "github.com/lejeunel/go-image-annotator/use-cases/label"
)

type Interactors struct {
	Label      *lbl.Interactors
	Collection *clc.Interactors
	Image      *im.Interactors
	Annotation *an.Interactors
}
