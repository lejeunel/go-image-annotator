package add_polygon

import (
	"time"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessAddPolygon(Response) {
	p.GotSuccess = true
}

type FakeAnnotationRepo struct {
	Err             error
	ErrOnAdd        bool
	GotImageId      im.ImageId
	GotCollectionId clc.CollectionId
	GotUserId       *u.UserId
	GotTime         *time.Time
	GotPolygon      a.Polygon
}

func (r *FakeAnnotationRepo) AddPolygon(imageId im.ImageId, collectionId clc.CollectionId, poly a.Polygon, userId *u.UserId, t *time.Time) error {
	if r.ErrOnAdd {
		return r.Err
	}
	r.GotImageId = imageId
	r.GotCollectionId = collectionId
	r.GotPolygon = poly
	r.GotUserId = userId
	r.GotTime = t
	return nil
}

type FakeLabelRepo struct {
	Err error
}

func (r *FakeLabelRepo) FindLabel(name string) (*lbl.Label, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	l := lbl.NewLabel(lbl.NewLabelId(), name)
	return &l, nil
}
