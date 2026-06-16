package update_label

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
	"time"
)

type FakeLabelRepo struct {
	Err                  error
	Returns              *lbl.Label
	FetchedLabelWithName string
}

func (r *FakeLabelRepo) FindLabel(name string) (*lbl.Label, error) {
	r.FetchedLabelWithName = name
	if r.Err != nil {
		return nil, r.Err
	}
	return r.Returns, nil
}

type FakeAnnotationRepo struct {
	Err                 error
	ErrOnFindLabel      bool
	ErrOnUpdate         bool
	UpdatedAnnotationId a.AnnotationId
	UpdatedLabelId      lbl.LabelId
	Returns             *lbl.Label
	GotUserId           *u.UserId
	GotTime             *time.Time
}

func (r *FakeAnnotationRepo) GroupOfAnnotation(a.AnnotationId) (*string, error) {
	group := "my-group"
	return &group, nil
}

func (r *FakeAnnotationRepo) UpdateLabelOfAnnotation(annotationId a.AnnotationId, labelId lbl.LabelId, userId *u.UserId, t *time.Time) error {
	if r.ErrOnUpdate {
		return r.Err
	}
	r.UpdatedAnnotationId = annotationId
	r.UpdatedLabelId = labelId
	r.GotUserId = userId
	r.GotTime = t
	return nil
}

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessUpdateLabel(Response) {
	p.GotSuccess = true
}
