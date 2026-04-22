package presenters

import (
	v "github.com/lejeunel/go-image-annotator-v2/app/annotator/view"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
)

type UpdateBoxPresenter struct {
	v.View
}

func (p UpdateBoxPresenter) SuccessUpdateBox(r updbox.Response) {
	p.View.UpdateBox(r)
}
func (p UpdateBoxPresenter) Error(err error) {}
