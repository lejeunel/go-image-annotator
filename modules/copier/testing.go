package copier

import (
	"time"

	"github.com/jonboulle/clockwork"
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func NewTestingCopier(opts ...Option) Copier {
	i := &Copier{
		Store:          &FakeImageStore{},
		CollectionRepo: &FakeCollectionRepo{},
		AnnotationRepo: &FakeAnnotationRepo{},
		clock:          clockwork.NewFakeClock(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

type FakeCollectionRepo struct {
	Err                    error
	ErrOnFindCollection    bool
	MissingCollection      bool
	CollectionWithoutGroup bool
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
	ErrOnAddPolygon       bool
	NumLabelsAdded        int
	NumBoundingboxesAdded int
	NumPolygonsAdded      int
}

type FakeImageStore struct {
	Err    error
	Return *im.Image
}

func (s *FakeImageStore) Find(baseImage im.BaseImage) (*im.Image, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	if s.Return != nil {
		return s.Return, nil
	}
	return &im.Image{}, nil
}

func (r *FakeCollectionRepo) FindCollectionByName(name string) (*clc.Collection, error) {
	if r.ErrOnFindCollection {
		return nil, r.Err
	}
	if r.MissingCollection {
		return nil, e.ErrNotFound
	}

	c := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	if r.CollectionWithoutGroup {
		return &c, nil
	}
	group := g.NewGroup(g.NewGroupId(), "my-group")
	c.Group = &group
	return &c, nil
}

func (r *FakeLabelRepo) FindLabel(name string) (*lbl.Label, error) {
	if r.ErrOnLabelExists {
		return nil, r.Err
	}
	if r.MissingLabel {
		return nil, e.ErrNotFound
	}
	l := lbl.NewLabel(lbl.NewLabelId(), name)
	return &l, nil
}

func (r *FakeAnnotationRepo) AddImageLabel(im.ImageId, clc.CollectionId, an.ImageLabel, *u.UserId, *time.Time) error {
	if r.ErrOnAddLabel {
		return r.Err
	}
	r.NumLabelsAdded += 1
	return nil
}

func (r *FakeAnnotationRepo) AddBoundingBox(im.ImageId, clc.CollectionId, an.BoundingBox, *u.UserId, *time.Time) error {
	if r.ErrOnAddBoundingBox {
		return r.Err
	}
	r.NumBoundingboxesAdded += 1
	return nil
}

func (r *FakeAnnotationRepo) AddPolygon(im.ImageId, clc.CollectionId, an.Polygon, *u.UserId, *time.Time) error {
	if r.ErrOnAddPolygon {
		return r.Err
	}
	r.NumPolygonsAdded += 1
	return nil
}
