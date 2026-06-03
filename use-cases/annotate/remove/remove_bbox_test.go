package remove

import (
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := NewInteractor(&FakeRepo{},
		WithAuth(auth.FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{Id: a.NewAnnotationId().String()},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingBoxShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrNotFound})
	itr.Execute(t.Context(), Request{Id: a.NewAnnotationId().String()}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestRemoveBox(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	itr := NewInteractor(repo)
	annotationId := a.NewAnnotationId()
	itr.Execute(t.Context(), Request{Id: annotationId.String()}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, annotationId, repo.Got)
}
