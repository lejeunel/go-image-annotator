package role

import (
	"testing"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewTestRoleRepo()
	repo.Db.Close()
	_, err := CreateRole(repo, "a-role")
	assert.ErrorIs(t, err, e.ErrInternal, "expected internal error")
}

func TestCreate(t *testing.T) {
	_, err := CreateRole(NewTestRoleRepo(), "a-role")
	assert.NoError(t, err, "expected no error on create but got")
}

func TestCreateWithMethods(t *testing.T) {
	repo := NewTestRoleRepo()
	name := "a-role"
	methods := []string{"a-method"}
	role := rl.NewRole(rl.NewRoleId(), name,
		rl.WithMethods(methods))
	err := repo.Create(role)
	assert.NoError(t, err)

	r, err := repo.Find(role.Name)
	assert.NoError(t, err)
	assert.Equal(t, methods, r.Methods)
}
