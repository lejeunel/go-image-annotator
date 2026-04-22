package presenters

import (
	v "github.com/lejeunel/go-image-annotator-v2/app/annotator/view"
	a "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
)

type AddBoxPresenter struct {
	View v.View
}

func (p AddBoxPresenter) SuccessAddBox(b a.BoundingBox) {
	p.View.AddBox(*v.MakeBoundingBox(&b, 0))
}
func (p AddBoxPresenter) Error(err error) {}
