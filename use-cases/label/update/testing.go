package update

import (
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
	"slices"
)

type FakeRepo struct {
	Names []string
	Got   Model
}

func (r *FakeRepo) Update(m Model) error {
	r.Got = m
	return nil
}

func (r *FakeRepo) Exists(n string) (bool, error) {
	if slices.Contains(r.Names, n) {
		return true, nil
	}
	return false, nil
}

type FakeErrRepo struct {
	err error
}

func (r *FakeErrRepo) Update(m Model) error {
	return r.err
}
func (r *FakeErrRepo) Exists(n string) (bool, error) {
	return false, r.err
}

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success(r Response) {
	p.GotSuccess = true
	p.Got = r
}

type FakeCollectionRepo struct {
	Err                 error
	ErrOnFindCollection bool
	MissingCollection   bool
}

type FakeLabelRepo struct {
	Err              error
	ErrOnLabelExists bool
	MissingLabel     bool
}

type FakeAnnotationRepo struct {
	Err                   error
	ErrOnAddBoundingBox   bool
	ErrOnAddLabel         bool
	NumLabelsAdded        int
	NumBoundingboxesAdded int
}

type FakeImageRepo struct {
	Err               error
	GotImage          bool
	ErrOnIngest       bool
	ErrOnFindHash     bool
	ErrOnDeleteImage  bool
	HashAlreadyExists bool
	NumDeletedImages  int
}

func (r *FakeCollectionRepo) FindCollectionByName(name string) (*clc.Collection, error) {
	if r.ErrOnFindCollection {
		return nil, r.Err
	}
	if r.MissingCollection {
		return nil, e.ErrNotFound
	}
	c := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	return c, nil
}

func (r *FakeLabelRepo) FindLabelByName(name string) (*lbl.Label, error) {
	if r.ErrOnLabelExists {
		return nil, r.Err
	}
	if r.MissingLabel {
		return nil, e.ErrNotFound
	}
	return lbl.NewLabel(lbl.NewLabelId(), name), nil
}

func (r *FakeImageRepo) FindImageByHash(hash string) (*im.Image, error) {
	if r.ErrOnFindHash {
		return nil, r.Err
	}
	if r.HashAlreadyExists {
		return im.NewImage(im.NewImageId(), *clc.NewCollection(clc.NewCollectionId(), "a-collection")), nil
	}
	return nil, e.ErrNotFound
}

func (r *FakeAnnotationRepo) AddImageLabel(an.AnnotationId, im.ImageId, clc.CollectionId, lbl.LabelId) error {
	if r.ErrOnAddLabel {
		return r.Err
	}
	r.NumLabelsAdded += 1
	return nil
}

func (r *FakeAnnotationRepo) AddBoundingBox(im.ImageId, clc.CollectionId, an.BoundingBox) error {
	if r.ErrOnAddBoundingBox {
		return r.Err
	}
	r.NumBoundingboxesAdded += 1
	return nil
}

func (r *FakeImageRepo) Delete(im.ImageId) error {
	if r.ErrOnDeleteImage {
		return r.Err
	}
	r.NumDeletedImages += 1
	return nil
}

func (r *FakeImageRepo) AddToCollection(im.ImageId, clc.CollectionId) error {
	if r.ErrOnIngest {
		return r.Err
	}
	return nil
}
