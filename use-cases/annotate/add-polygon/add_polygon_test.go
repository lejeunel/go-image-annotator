package add_polygon

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
	"github.com/stretchr/testify/assert"
)

func CreateImage() im.Image {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	return im.NewImage(im.NewImageId(), collection)
}

func TestHandleAuthError(t *testing.T) {
	image := CreateImage()
	group := g.NewGroup(g.NewGroupId(), "my-group")
	image.Collection.Group = &group
	itr := New(&fk.ImageStore{Return: &image},
		&fk.AnnotationRepo{},
		&fk.LabelRepo{},
		WithAuth(auth.FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: im.NewImageId().String()},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestErrOnImageRetrievalShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageStore{Err: e.ErrInternal},
		&fk.AnnotationRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestErrOnFindLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageStore{},
		&fk.AnnotationRepo{},
		&fk.LabelRepo{ErrOnFind: e.ErrInternal},
	)
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func CreateTestAddPolygonRequest() Request {
	return Request{ImageId: im.NewImageId().String(), Collection: "a-collection",

		Label: "a-label", Points: a.Points{Coordinates: [][2]float32{{0, 0}, {1, 1}}}}
}

func TestErrOnAddPolygonShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageStore{},
		&fk.AnnotationRepo{ErrOnAddPoly: e.ErrInternal},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), CreateTestAddPolygonRequest(), p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestAddUserIdFromContext(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.AnnotationRepo{}
	itr := New(&fk.ImageStore{}, repo,
		&fk.LabelRepo{})
	user := u.NewUser("user@example.com")
	ctx := u.AppendUserToContext(t.Context(), user)
	itr.Execute(ctx, CreateTestAddPolygonRequest(), p)
	assert.NotNil(t, repo.GotUserId)
	assert.Equal(t, user.Id, *repo.GotUserId)
}

func TestTime(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.AnnotationRepo{}
	now := time.Now()
	itr := New(&fk.ImageStore{},
		repo,
		&fk.LabelRepo{},
		WithClock(clockwork.NewFakeClockAt(now)))
	itr.Execute(t.Context(), CreateTestAddPolygonRequest(), p)
	assert.NotNil(t, repo.GotTime)
	assert.Equal(t, now, *repo.GotTime)
}

func TestAddPolygon(t *testing.T) {
	p := &FakePresenter{}
	repo := fk.AnnotationRepo{}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	image := im.NewImage(im.NewImageId(), collection)
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	req := CreateTestAddPolygonRequest()
	req.ImageId = image.Id.String()
	itr := New(&fk.ImageStore{Return: &image},
		&repo,
		&fk.LabelRepo{Return: label},
	)
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, req.ImageId, repo.GotImageId.String())
	assert.Equal(t, collection.Id, repo.GotCollectionId)
	assert.Equal(t, req.Label, repo.GotPolygon.Label.Name)
	assert.Equal(t, req.Points, repo.GotPolygon.Points)

}
