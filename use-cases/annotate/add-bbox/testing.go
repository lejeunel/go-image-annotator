package add_bbox

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessAddBox(a.BoundingBox) {
	p.GotSuccess = true
}

type FakeRepo struct {
	Err             error
	ErrOnAdd        bool
	ErrOnFindLabel  bool
	GotImageId      im.ImageId
	GotCollectionId clc.CollectionId
	GotBox          a.BoundingBox
}

func (r *FakeRepo) AddBoundingBox(imageId im.ImageId, collectionId clc.CollectionId, box a.BoundingBox) error {
	if r.ErrOnAdd {
		return r.Err
	}
	r.GotImageId = imageId
	r.GotCollectionId = collectionId
	r.GotBox = box
	return nil
}

func (r *FakeRepo) FindLabelByName(name string) (*lbl.Label, error) {
	if r.ErrOnFindLabel {
		return nil, r.Err
	}
	return lbl.NewLabel(lbl.NewLabelId(), name), nil
}
