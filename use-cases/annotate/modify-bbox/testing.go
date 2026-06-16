package modify_bbox

import (
	"time"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err            error
	ErrOnUpdate    bool
	ErrOnFindLabel bool
	Got            a.BoundingBoxUpdatables
	GotUserId      *u.UserId
	GotTime        *time.Time
	Label          lbl.Label
	NoGroup        bool
}

func (r *FakeRepo) UpdateBoundingBox(id a.AnnotationId, u a.BoundingBoxUpdatables, userId *u.UserId, t *time.Time) error {
	if r.ErrOnUpdate {
		return r.Err
	}
	r.Got = u
	r.GotUserId = userId
	r.GotTime = t
	return nil
}

func (r *FakeRepo) FindLabel(name string) (*lbl.Label, error) {
	if r.ErrOnFindLabel {
		return nil, r.Err
	}
	return &r.Label, nil
}
func (r *FakeRepo) GroupOfAnnotation(id a.AnnotationId) (*string, error) {
	if r.NoGroup {
		return nil, nil
	}
	group := "my-group"
	return &group, nil
}

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessUpdateBox(Response) {
	p.GotSuccess = true
}
