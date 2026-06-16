package assign_label

import (
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
	"time"
)

type FakeRepo struct {
	Err            error
	ErrOnAddLabel  bool
	ErrOnFindLabel bool
	MissingLabel   bool

	AddedLabelId        lbl.LabelId
	AddedOnImageId      im.ImageId
	AddedOnCollectionId clc.CollectionId

	ReturnLabel lbl.Label
	GotUserId   *u.UserId
	GotTime     *time.Time
}

func (r *FakeRepo) FindLabel(string) (*lbl.Label, error) {
	if r.MissingLabel {
		return nil, e.ErrNotFound
	}
	if r.ErrOnFindLabel {
		return nil, r.Err
	}
	return &r.ReturnLabel, nil
}

// AddImageLabel(im.ImageId, clc.CollectionId, an.ImageLabel, *u.UserId, *time.Time) error
func (r *FakeRepo) AddImageLabel(imageId im.ImageId, collectionId clc.CollectionId, imageLabel an.ImageLabel, userId *u.UserId, t *time.Time) error {
	if r.ErrOnAddLabel {
		return r.Err
	}
	r.AddedLabelId = imageLabel.Label.Id
	r.AddedOnImageId = imageId
	r.AddedOnCollectionId = collectionId
	r.GotUserId = userId
	r.GotTime = t
	return nil

}

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessAddLabel(r Response) {
	p.Got = r
	p.GotSuccess = true
}
