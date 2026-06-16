package remove

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err     error
	Got     a.AnnotationId
	NoGroup bool
}

func (r *FakeRepo) GroupOfAnnotation(annotationId a.AnnotationId) (*string, error) {
	if r.NoGroup {
		return nil, nil
	}
	group := "my-group"
	return &group, nil
}

func (r *FakeRepo) RemoveAnnotation(annotationId a.AnnotationId) error {
	if r.Err != nil {
		return r.Err
	}
	r.Got = annotationId
	return nil
}

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessDeleteAnnotation(Response) {
	p.GotSuccess = true
}
