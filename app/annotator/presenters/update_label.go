package presenters

import (
	v "github.com/lejeunel/go-image-annotator-v2/app/annotator/view"
	updlbl "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/update-label"
)

type UpdateLabelOfAnnotationPresenter struct {
	v.View
}

func (p UpdateLabelOfAnnotationPresenter) SuccessUpdateLabel(r updlbl.Response) {
}
func (p UpdateLabelOfAnnotationPresenter) Error(err error) {}
