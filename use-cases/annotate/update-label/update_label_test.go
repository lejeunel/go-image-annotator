package update_label

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
	"github.com/stretchr/testify/assert"
)

func CreateTestRequest() Request {
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	return Request{AnnotationId: a.NewAnnotationId().String(), Label: newLabel.Name}
}

func TestHandleAuthError(t *testing.T) {
	lbl := lbl.NewLabel(lbl.NewLabelId(), "my-label")
	itr := New(&FakeRepo{Returns: &lbl},
		WithAuth(auth.FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), CreateTestRequest(), p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrOnFindLabel(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{Err: e.ErrNotFound, ErrOnFindLabel: true})
	itr.Execute(t.Context(), CreateTestRequest(), p)
	assert.False(t, p.GotSuccess)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
}

func TestHandleErrOnUpdateLabel(t *testing.T) {
	p := &FakePresenter{}
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	itr := New(&FakeRepo{Returns: &label, Err: e.ErrNotFound, ErrOnUpdate: true})
	itr.Execute(t.Context(), CreateTestRequest(), p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestFetchLabelFromName(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeRepo{Returns: &newLabel}
	itr := New(repo)
	itr.Execute(t.Context(), CreateTestRequest(), p)
	assert.Equal(t, repo.FetchedLabelWithName, newLabel.Name, "label name")
}
func TestAddUserIdFromContext(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeRepo{Returns: &newLabel}
	itr := New(repo)
	user := u.NewUser("user@example.com")
	ctx := context.WithValue(t.Context(), u.UserContextKey, &user)
	itr.Execute(ctx, CreateTestRequest(), p)
	assert.NotNil(t, repo.GotUserId)
	assert.Equal(t, user.Id, *repo.GotUserId)
}

func TestTime(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeRepo{Returns: &newLabel}
	now := time.Now()
	itr := New(repo, WithClock(clockwork.NewFakeClockAt(now)))
	itr.Execute(t.Context(), CreateTestRequest(), p)
	assert.NotNil(t, repo.GotTime)
	assert.Equal(t, now, *repo.GotTime)
}

func TestUpdateLabel(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeRepo{Returns: &newLabel}
	itr := New(repo)
	req := CreateTestRequest()
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, req.AnnotationId, repo.UpdatedAnnotationId.String())
	assert.Equal(t, repo.UpdatedLabelId, newLabel.Id)
}
