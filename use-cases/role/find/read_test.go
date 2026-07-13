package find

import (
	"github.com/stretchr/testify/assert"
	"testing"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestRead(t *testing.T) {
	role := rl.NewRole(rl.NewRoleId(), "my-role")
	repo := &FakeRepo{Role: role}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), Request{Name: role.Name}, p)
	assert.Equal(t, Response{Name: role.Name}, p.Got)
}

func TestReadNonExistingShouldFail(t *testing.T) {
	role := rl.NewRole(rl.NewRoleId(), "my-role")
	repo := &FakeRepo{Role: role}
	p := &FakePresenter{}
	itr := New(repo)
	req := Request{Name: "non-existing-role"}
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
}
