package delete

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	itr := NewInteractor(&FakeRepo{}, WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(st.FakeProvider{}, Request{}, p)
	assert.Equal(t, true, p.GotAuthErr, "auth error")
	assert.Equal(t, false, p.GotSuccess)
}

func TestDeleteNonExistingCollectionShouldFail(t *testing.T) {

	name := "my-collection"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Missing: true})
	itr.Execute(st.FakeProvider{}, Request{Name: name}, p)
	assert.Equal(t, p.GotNotFoundErr, true, "collection not found")
	assert.Equal(t, p.GotSuccess, false)
}

func TestDeleteCollectionWithAssociatedResourcesShouldFail(t *testing.T) {

	name := "my-collection"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{IsPopulated_: true})
	itr.Execute(st.FakeProvider{}, Request{Name: name}, p)
	assert.Equal(t, p.GotDependencyErr, true, "expected dependency error, but got none")
	assert.Equal(t, p.GotSuccess, false)
}

func TestHandleInternalErrorOnDelete(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnDelete: true, Err: e.ErrInternal})
	itr.Execute(st.FakeProvider{}, Request{}, p)
	assert.Equal(t, p.GotInternalErr, true, "expected internal error, but got none")
}

func TestDeleteCollection(t *testing.T) {

	name := "my-collection"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{})
	itr.Execute(st.FakeProvider{}, Request{Name: name}, p)
	assert.Equal(t, p.GotSuccess, true, "expected success, but did not")
}
