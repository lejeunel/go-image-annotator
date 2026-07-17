package delete

import (
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	group := "a-group"
	itr := New(&fk.CollectionRepo{}, &fk.GroupRepo{Return: group},
		WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{}, &fk.GroupRepo{})
	itr.Execute(t.Context(), "my-collection", p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteCollectionWithAssociatedResourcesShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{ExistingNames: []string{"my-collection"},
		IsPopulated_: true}, &fk.GroupRepo{})
	itr.Execute(t.Context(), "my-collection", p)
	assert.True(t, p.GotDependencyErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrorOnDelete(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{ExistingNames: []string{"my-collection"},
		ErrOnDelete: e.ErrInternal}, &fk.GroupRepo{})
	itr.Execute(t.Context(), "my-collection", p)
	assert.True(t, p.GotInternalErr)
}

func TestDeleteCollection(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{ExistingNames: []string{"my-collection"}},
		&fk.GroupRepo{ErrOnGetGroupOfCollection: e.ErrNotFound})
	itr.Execute(t.Context(), "my-collection", p)
	assert.True(t, p.GotSuccess)
}
