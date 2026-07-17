package delete

import (
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.RoleRepo{}, WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteNonExistingRoleShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{})
	itr.Execute(t.Context(), "my-role", p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteRoleAssignedToUsersShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{ExistingNames: []string{"my-role"}, IsAssigned_: true})
	itr.Execute(t.Context(), "my-role", p)
	assert.True(t, p.GotDependencyErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrorOnDelete(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{ExistingNames: []string{"my-role"}, ErrOnDelete: e.ErrInternal})
	itr.Execute(t.Context(), "my-role", p)
	assert.True(t, p.GotInternalErr)
}

func TestDeleteRole(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{ExistingNames: []string{"my-role"}})
	itr.Execute(t.Context(), "my-role", p)
	assert.True(t, p.GotSuccess)
}
