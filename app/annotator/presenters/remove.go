package presenters

import (
	v "github.com/lejeunel/go-image-annotator-v2/app/annotator/view"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
)

type RemoveAnnotationPresenter struct {
	v.View
}

func (p RemoveAnnotationPresenter) SuccessDeleteAnnotation(r del.Response) {
	p.View.DeleteAnnotation(r)
}
func (p RemoveAnnotationPresenter) Error(err error) {}
