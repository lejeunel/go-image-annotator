package list

import (
	r "github.com/lejeunel/go-image-annotator/entities/role"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{ErrOnList: e.ErrInternal})
	itr.Execute(t.Context(), p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestList(t *testing.T) {
	r0 := r.NewRole(r.NewRoleId(), "a-role")
	r1 := r.NewRole(r.NewRoleId(), "another-role")
	repo := &fk.RoleRepo{ReturnList: []r.Role{r0, r1}}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), p)
	assert.Equal(t, 2, len(p.Got))
}
