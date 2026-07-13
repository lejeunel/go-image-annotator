package delete

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&FakeRepo{}, WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteNonExistingRoleShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{Missing: true})
	itr.Execute(t.Context(), Request{Name: "my-role"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteRoleAssignedToUsersShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{IsAssigned_: true})
	itr.Execute(t.Context(), Request{Name: "my-role"}, p)
	assert.True(t, p.GotDependencyErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrorOnDelete(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{ErrOnDelete: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
}

func TestDeleteRole(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{})
	itr.Execute(t.Context(), Request{Name: "my-role"}, p)
	assert.True(t, p.GotSuccess)
}
