package update_label

import (
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
	itr := New(&FakeAnnotationRepo{Returns: &lbl},
		&FakeLabelRepo{},
		WithAuth(auth.FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), CreateTestRequest(), p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrOnFindLabel(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeAnnotationRepo{}, &FakeLabelRepo{Err: e.ErrNotFound})
	itr.Execute(t.Context(), CreateTestRequest(), p)
	assert.False(t, p.GotSuccess)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
}

func TestHandleErrOnUpdateLabel(t *testing.T) {
	p := &FakePresenter{}
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	itr := New(&FakeAnnotationRepo{Err: e.ErrNotFound, ErrOnUpdate: true},
		&FakeLabelRepo{Returns: &label})
	itr.Execute(t.Context(), CreateTestRequest(), p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestAddUserIdFromContext(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeAnnotationRepo{}
	itr := New(repo, &FakeLabelRepo{Returns: &newLabel})
	user := u.NewUser("user@example.com")
	ctx := u.AppendUserToContext(t.Context(), user)
	itr.Execute(ctx, CreateTestRequest(), p)
	assert.NotNil(t, repo.GotUserId)
	assert.Equal(t, user.Id, *repo.GotUserId)
}

func TestTime(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeAnnotationRepo{}
	now := time.Now()
	itr := New(repo, &FakeLabelRepo{Returns: &newLabel},
		WithClock(clockwork.NewFakeClockAt(now)))
	itr.Execute(t.Context(), CreateTestRequest(), p)
	assert.NotNil(t, repo.GotTime)
	assert.Equal(t, now, *repo.GotTime)
}

func TestUpdateLabelNoGroup(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeAnnotationRepo{NoGroup: true}
	itr := New(repo, &FakeLabelRepo{Returns: &newLabel})
	req := CreateTestRequest()
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, req.AnnotationId, repo.UpdatedAnnotationId.String())
	assert.Equal(t, repo.UpdatedLabelId, newLabel.Id)
}

func TestUpdateLabel(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeAnnotationRepo{}
	itr := New(repo, &FakeLabelRepo{Returns: &newLabel})
	req := CreateTestRequest()
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, req.AnnotationId, repo.UpdatedAnnotationId.String())
	assert.Equal(t, repo.UpdatedLabelId, newLabel.Id)
}
