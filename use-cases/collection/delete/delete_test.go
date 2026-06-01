package delete

import (
	"context"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	itr := NewInteractor(&FakeRepo{}, WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(context.Background(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Missing: true})
	itr.Execute(context.Background(), Request{Name: "my-collection"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteCollectionWithAssociatedResourcesShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{IsPopulated_: true})
	itr.Execute(context.Background(), Request{Name: "my-collection"}, p)
	assert.True(t, p.GotDependencyErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrorOnDelete(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnDelete: true, Err: e.ErrInternal})
	itr.Execute(context.Background(), Request{}, p)
	assert.True(t, p.GotInternalErr)
}

func TestDeleteCollection(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{})
	itr.Execute(context.Background(), Request{Name: "my-collection"}, p)
	assert.True(t, p.GotSuccess)
}
