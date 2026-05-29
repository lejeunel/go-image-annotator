package unassign_label

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err              error
	ErrOnRemoveLabel bool
	RemovedLabel     bool
}

func (r *FakeRepo) RemoveImageLabel(im.ImageId, clc.CollectionId, lbl.LabelId) error {
	if r.ErrOnRemoveLabel {
		return r.Err
	}
	r.RemovedLabel = true
	return nil

}

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success(Response) {
	p.GotSuccess = true
}
