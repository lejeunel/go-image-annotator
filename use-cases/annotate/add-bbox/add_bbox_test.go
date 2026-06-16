package add_bbox

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	st "github.com/lejeunel/go-image-annotator/modules/image-store"
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
	itr := New(&st.FakeImageStore{Return: &image},
		&FakeAnnotationRepo{},
		&FakeLabelRepo{},
		WithAuth(auth.FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestErrOnImageRetrievalShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&st.FakeImageStore{Err: e.ErrInternal},
		&FakeAnnotationRepo{},
		&FakeLabelRepo{})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestErrOnFindLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&st.FakeImageStore{},
		&FakeAnnotationRepo{},
		&FakeLabelRepo{Err: e.ErrInternal},
	)
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func CreateTestAddBoxRequest() Request {
	return Request{ImageId: im.NewImageId().String(), Collection: "a-collection",
		Label: "a-label", Xc: float32(1.0), Yc: float32(1.0), Width: float32(3.0),
		Height: float32(3.0), Angle: float32(32)}
}

func TestValidationErrShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&st.FakeImageStore{}, &FakeAnnotationRepo{},
		&FakeLabelRepo{})
	req := CreateTestAddBoxRequest()
	req.Width = -999
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnAddBoxShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&st.FakeImageStore{},
		&FakeAnnotationRepo{ErrOnAdd: true, Err: e.ErrInternal},
		&FakeLabelRepo{})
	itr.Execute(t.Context(), CreateTestAddBoxRequest(), p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestAddUserIdFromContext(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeAnnotationRepo{}
	itr := New(&st.FakeImageStore{}, repo,
		&FakeLabelRepo{})
	user := u.NewUser("user@example.com")
	ctx := context.WithValue(t.Context(), u.UserContextKey, &user)
	itr.Execute(ctx, CreateTestAddBoxRequest(), p)
	assert.NotNil(t, repo.GotUserId)
	assert.Equal(t, user.Id, *repo.GotUserId)
}

func TestTime(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeAnnotationRepo{}
	now := time.Now()
	itr := New(&st.FakeImageStore{},
		repo,
		&FakeLabelRepo{},
		WithClock(clockwork.NewFakeClockAt(now)))
	itr.Execute(t.Context(), CreateTestAddBoxRequest(), p)
	assert.NotNil(t, repo.GotTime)
	assert.Equal(t, now, *repo.GotTime)
}

func TestAddBoundingBox(t *testing.T) {
	p := &FakePresenter{}
	repo := FakeAnnotationRepo{}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	image := im.NewImage(im.NewImageId(), collection)
	req := Request{ImageId: image.Id.String(), Collection: collection.Name,
		Label: "a-label", Xc: float32(1.0), Yc: float32(1.0), Width: float32(3.0),
		Height: float32(3.0), Angle: float32(32)}

	itr := New(&st.FakeImageStore{Return: &image},
		&repo,
		&FakeLabelRepo{},
	)
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, req.ImageId, repo.GotImageId.String())
	assert.Equal(t, collection.Id, repo.GotCollectionId)
	assert.Equal(t, req.Label, repo.GotBox.Label.Name)
	assert.Equal(t, req.Xc, repo.GotBox.Xc)
	assert.Equal(t, req.Yc, repo.GotBox.Yc)
	assert.Equal(t, req.Width, repo.GotBox.Width)
	assert.Equal(t, req.Height, repo.GotBox.Height)
	assert.Equal(t, req.Angle, repo.GotBox.Angle)

}
