package clone

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	el "github.com/lejeunel/go-image-annotator/modules/event-logger"
	st "github.com/lejeunel/go-image-annotator/modules/image-store"
	testing "github.com/lejeunel/go-image-annotator/shared/testing"
	fk "github.com/lejeunel/go-image-annotator/use-cases/fakes"
	"iter"
	"time"
)

type FakeImageRepo struct{}

func (r *FakeImageRepo) AddToCollection(im.ImageId, clc.CollectionId) error {
	return nil
}

func (r *FakeImageRepo) Iterate(im.FilteringParams, int) iter.Seq2[im.BaseImage, error] {
	return nil
}

type FakeCollectionRepo struct{}

func (r *FakeCollectionRepo) Create(clc.Collection) error {
	return nil
}
func (r *FakeCollectionRepo) Exists(string) (bool, error) {
	return false, nil
}

type FakeAnnotationRepo struct{}

func (r *FakeAnnotationRepo) AddImageLabel(im.ImageId, clc.CollectionId, a.ImageLabel, *u.UserId, *time.Time) error {
	return nil
}
func (r *FakeAnnotationRepo) AddBoundingBox(im.ImageId, clc.CollectionId, a.BoundingBox, *u.UserId, *time.Time) error {
	return nil
}
func (r *FakeAnnotationRepo) AddPolygon(im.ImageId, clc.CollectionId, a.Polygon, *u.UserId, *time.Time) error {
	return nil
}

type FakeImageStore struct{}
type FakeEventLogger struct{}

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	testing.TestingErrPresenter
}

func (p *FakePresenter) SuccessSubmitCloneTask(r Response) {
	p.Got = r
	p.GotSuccess = true
}

func NewTestingCloner() Interactor {
	return New(&FakeImageRepo{}, &FakeCollectionRepo{}, &FakeAnnotationRepo{}, &fk.GroupRepo{}, &st.FakeImageStore{}, &el.FakeEventLogger{})
}
