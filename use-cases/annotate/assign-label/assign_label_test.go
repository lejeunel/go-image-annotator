package assign_label

import (
	"testing"
	"time"

	"context"

	"github.com/jonboulle/clockwork"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
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
	itr := New(&FakeAnnotationRepo{},
		&FakeLabelRepo{},
		&st.FakeImageStore{Return: &image},
		WithAuth(auth.FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleNotFoundErrOnImageRetrieval(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeAnnotationRepo{},
		&FakeLabelRepo{},
		&st.FakeImageStore{Err: e.ErrNotFound})
	itr.Execute(t.Context(), Request{im.NewImageId().String(), "a-collection", "a-label"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnImageRetrieval(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeAnnotationRepo{},
		&FakeLabelRepo{},
		&st.FakeImageStore{Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{im.NewImageId().String(), "a-collection", "a-label"}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestAssignNonExistingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	image := CreateImage()
	itr := New(&FakeAnnotationRepo{},
		&FakeLabelRepo{Err: e.ErrNotFound},
		&st.FakeImageStore{Return: &image})
	itr.Execute(t.Context(), Request{image.Id.String(), image.Collection.Name, "a-label"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}
func TestAddUserIdFromContext(t *testing.T) {
	p := &FakePresenter{}
	image := CreateImage()
	repo := &FakeAnnotationRepo{}
	itr := New(repo, &FakeLabelRepo{}, &st.FakeImageStore{Return: &image})
	user := u.NewUser("user@example.com")
	ctx := context.WithValue(t.Context(), u.UserContextKey, &user)
	itr.Execute(ctx, Request{im.NewImageId().String(), "a-collection", "a-label"}, p)
	assert.NotNil(t, repo.GotUserId)
	assert.Equal(t, user.Id, *repo.GotUserId)
}
func TestTime(t *testing.T) {
	p := &FakePresenter{}
	image := CreateImage()
	repo := &FakeAnnotationRepo{}
	now := time.Now()
	itr := New(repo, &FakeLabelRepo{},
		&st.FakeImageStore{Return: &image}, WithClock(clockwork.NewFakeClockAt(now)))
	itr.Execute(t.Context(), Request{im.NewImageId().String(), "a-collection", "a-label"}, p)
	assert.NotNil(t, repo.GotTime)
	assert.Equal(t, now, *repo.GotTime)
}

func TestAssignLabelToImage(t *testing.T) {
	p := &FakePresenter{}
	image := CreateImage()
	label := lbl.NewLabel(lbl.NewLabelId(), "al-label")
	req := Request{ImageId: image.Id.String(),
		Collection: image.Collection.Name,
		Label:      label.Name}
	repo := &FakeAnnotationRepo{}
	itr := New(repo,
		&FakeLabelRepo{ReturnLabel: label},
		&st.FakeImageStore{Return: &image})
	itr.Execute(t.Context(), req, p)
	resp := p.Got
	assert.Equal(t, resp.Label, req.Label, "label")
	assert.Equal(t, resp.Collection, req.Collection, "collection")
	assert.Equal(t, resp.ImageId, req.ImageId, "image id")
	assert.Equal(t, repo.AddedLabelId, label.Id, "added label id")
	assert.Equal(t, repo.AddedOnImageId, image.Id, "added on image id")
	assert.Equal(t, repo.AddedOnCollectionId, image.Collection.Id, "added on collection id")
}
