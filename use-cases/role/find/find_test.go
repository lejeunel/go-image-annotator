package find

import (
	"testing"

	"github.com/stretchr/testify/assert"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{ErrOnFind: e.ErrInternal})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotInternalErr)
}

func TestRead(t *testing.T) {
	methods := []string{"a-method"}
	role := rl.NewRole(rl.NewRoleId(), "my-role", rl.WithMethods(methods))
	repo := &fk.RoleRepo{ExistingRoles: []rl.Role{role}}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), role.Name, p)
	assert.Equal(t, role.Name, p.Got.Name)
	assert.Equal(t, role.Methods, p.Got.Methods)
}
