package modify_bbox

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
	"github.com/stretchr/testify/assert"
)

func CreateRequestAndUpdatable() (Request, a.BoundingBoxUpdatables, lbl.Label) {
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	req := Request{AnnotationId: a.NewAnnotationId().String(), Xc: 1, Yc: 1, Width: 1, Height: 1,
		Angle: -1, Label: label.Name}
	upd := a.BoundingBoxUpdatables{LabelId: label.Id, Xc: req.Xc,
		Yc: req.Yc, Width: req.Width, Height: req.Height, Angle: req.Angle}
	return req, upd, label
}

func AssertUpdated(t *testing.T, expected, got a.BoundingBoxUpdatables) {
	assert.Equal(t, expected.Xc, got.Xc)
	assert.Equal(t, expected.Yc, got.Yc)
	assert.Equal(t, expected.Width, got.Width)
	assert.Equal(t, expected.Height, got.Height)
	assert.Equal(t, expected.Angle, got.Angle)
	assert.Equal(t, expected.LabelId, got.LabelId)
}

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.AnnotationRepo{},
		&fk.LabelRepo{},
		WithAuth(auth.FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{AnnotationId: a.NewAnnotationId().String()},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.AnnotationRepo{},
		&fk.LabelRepo{ErrOnFind: e.ErrNotFound})
	itr.Execute(t.Context(),
		Request{AnnotationId: a.NewAnnotationId().String()},
		p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnFindLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.AnnotationRepo{}, &fk.LabelRepo{ErrOnFind: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestValidationErrShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.AnnotationRepo{}, &fk.LabelRepo{})
	itr.Execute(t.Context(),
		Request{AnnotationId: a.NewAnnotationId().String(), Xc: 1, Yc: 1, Width: -999, Height: 1}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}
func TestNotFoundErrOnUpdateShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.AnnotationRepo{ErrOnUpdate: e.ErrNotFound},
		&fk.LabelRepo{})
	req, _, _ := CreateRequestAndUpdatable()
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateWithDefaultGroup(t *testing.T) {
	p := &FakePresenter{}
	req, upd, label := CreateRequestAndUpdatable()
	repo := &fk.AnnotationRepo{NoGroup: true}
	itr := New(repo, &fk.LabelRepo{Return: label})
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotSuccess)
	AssertUpdated(t, upd, repo.GotUpdatableBox)
}

func TestUpdateWithUserIdFromContext(t *testing.T) {
	p := &FakePresenter{}
	req, _, label := CreateRequestAndUpdatable()
	repo := &fk.AnnotationRepo{NoGroup: true}
	itr := New(repo, &fk.LabelRepo{Return: label})
	user := u.NewUser("user@example.com")
	ctx := u.AppendUserToContext(t.Context(), user)
	itr.Execute(ctx, req, p)
	assert.NotNil(t, repo.GotUserId)
	assert.Equal(t, user.Id, *repo.GotUserId)
}

func TestTime(t *testing.T) {
	p := &FakePresenter{}
	req, _, label := CreateRequestAndUpdatable()
	repo := &fk.AnnotationRepo{NoGroup: true}
	now := time.Now()
	itr := New(repo, &fk.LabelRepo{Return: label},
		WithClock(clockwork.NewFakeClockAt(now)))
	itr.Execute(t.Context(), req, p)
	assert.NotNil(t, repo.GotTime)
	assert.Equal(t, now, *repo.GotTime)
}

func TestUpdate(t *testing.T) {
	p := &FakePresenter{}
	req, upd, label := CreateRequestAndUpdatable()
	repo := &fk.AnnotationRepo{}
	itr := New(repo, &fk.LabelRepo{Return: label})
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotSuccess)
	AssertUpdated(t, upd, repo.GotUpdatableBox)
}
