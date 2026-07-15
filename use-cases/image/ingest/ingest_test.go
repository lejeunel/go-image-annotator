package ingest

import (
	"testing"

	ig "github.com/lejeunel/go-image-annotator/modules/ingester"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := NewTestingInteractor(&FakeCollectionRepo{}, WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), ig.Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor(&FakeCollectionRepo{MissingCollection: true})
	itr.Execute(t.Context(), ig.Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrorOnCollectionExistsCheck(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor(&FakeCollectionRepo{ErrOnFindCollection: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), ig.Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}
