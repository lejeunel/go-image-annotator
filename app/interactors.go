package app

import (
	an "github.com/lejeunel/go-image-annotator/use-cases/annotate"
	clc "github.com/lejeunel/go-image-annotator/use-cases/collection"
	grp "github.com/lejeunel/go-image-annotator/use-cases/group"
	im "github.com/lejeunel/go-image-annotator/use-cases/image"
	lbl "github.com/lejeunel/go-image-annotator/use-cases/label"
	usr "github.com/lejeunel/go-image-annotator/use-cases/user"
)

type Interactors struct {
	Label      lbl.Interactors
	Collection clc.Interactors
	Image      im.Interactors
	Annotation an.Interactors
	Group      grp.Interactors
	User       usr.Interactors
}
