package delete

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	group := "a-group"
	itr := New(&FakeCollectionRepo{}, &FakeGroupRepo{Return: &group},
		WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeCollectionRepo{Missing: true}, &FakeGroupRepo{})
	itr.Execute(t.Context(), Request{Name: "my-collection"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteCollectionWithAssociatedResourcesShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeCollectionRepo{IsPopulated_: true}, &FakeGroupRepo{})
	itr.Execute(t.Context(), Request{Name: "my-collection"}, p)
	assert.True(t, p.GotDependencyErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrorOnDelete(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeCollectionRepo{ErrOnDelete: true, Err: e.ErrInternal}, &FakeGroupRepo{})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
}

func TestDeleteCollection(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeCollectionRepo{}, &FakeGroupRepo{})
	itr.Execute(t.Context(), Request{Name: "my-collection"}, p)
	assert.True(t, p.GotSuccess)
}
