package update_label

import (
	a "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator-v2/entities/label"
	t "github.com/lejeunel/go-image-annotator-v2/shared/testing"
)

type FakeRepo struct {
	Err                  error
	ErrOnFindLabel       bool
	ErrOnUpdate          bool
	FetchedLabelWithName string
	UpdatedAnnotationId  a.AnnotationId
	UpdatedLabelId       lbl.LabelId
	Returns              *lbl.Label
}

func (r *FakeRepo) FindLabelByName(name string) (*lbl.Label, error) {
	r.FetchedLabelWithName = name
	if r.ErrOnFindLabel {
		return nil, r.Err
	}
	return r.Returns, nil
}
func (r *FakeRepo) UpdateLabelOfAnnotation(annotationId a.AnnotationId, labelId lbl.LabelId) error {
	r.UpdatedAnnotationId = annotationId
	r.UpdatedLabelId = labelId
	if r.ErrOnUpdate {
		return r.Err
	}
	return nil
}

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessUpdateLabel(Response) {
	p.GotSuccess = true
}
