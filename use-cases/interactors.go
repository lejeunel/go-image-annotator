package interactors

import (
	an "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate"
	clc "github.com/lejeunel/go-image-annotator-v2/use-cases/collection"
	im "github.com/lejeunel/go-image-annotator-v2/use-cases/image"
	lbl "github.com/lejeunel/go-image-annotator-v2/use-cases/label"
)

type Interactors struct {
	Label      *lbl.Interactors
	Collection *clc.Interactors
	Image      *im.Interactors
	Annotation *an.Interactors
}
