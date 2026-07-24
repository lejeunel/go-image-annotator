package find

import (
	"github.com/stretchr/testify/assert"
	"testing"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestRead(t *testing.T) {
	role := rl.NewRole(rl.NewRoleId(), "my-role")
	repo := &fk.RoleRepo{Return: role}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), role.Name, p)
	assert.Equal(t, role.Name, p.Got.Name)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{ErrOnFind: e.ErrInternal})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotInternalErr)
}
